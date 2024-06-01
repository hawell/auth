package server

import (
	"auth/database"
	"auth/mailer"
	"auth/recaptcha"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

var (
	serverConfig = Config{
		BindAddress:   "localhost:8080",
		ReadTimeout:   100,
		WriteTimeout:  100,
		MaxBodyBytes:  10000000,
		WebServer:     "z42.com",
		HtmlTemplates: "../templates/*.tmpl",
		Recaptcha: &recaptcha.Config{
			SecretKey: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			Server:    "http://127.0.0.1:9798",
			Bypass:    false,
		},
	}
	connectionStr = "admin:admin@tcp(127.0.0.1:3306)/auth"
	db            *database.Database
	client        *http.Client
	tokens        map[database.ObjectId]string
	users         = []database.User{
		{
			Email:    "apiUser1",
			Password: "apiUser1",
		},
		{
			Email:    "apiUser2",
			Password: "apiUser2",
		},
		{
			Email:    "apiUser3",
			Password: "apiUser3",
		},
	}
	recaptchaServer = recaptcha.NewMockServer("127.0.0.1:9798")
)

func TestSignup(t *testing.T) {
	initialize(t)

	// add new user
	body := `{"email": "user1@example.com", "password": "password", "recaptcha_token": "123456"}`
	path := "/auth/signup"
	resp := execRequest("", http.MethodPost, path, body)
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))
	err := resp.Body.Close()
	Expect(err).To(BeNil())

	// check new user status is pending
	user, err := db.GetUser("user1@example.com")
	Expect(err).To(BeNil())
	Expect(user.Email).To(Equal("user1@example.com"))
	Expect(user.Status).To(Equal(database.UserStatusPending))

}

func TestVerify(t *testing.T) {
	initialize(t)

	_, code, err := db.AddUser(database.NewUser{
		Email:    "user2@example.com",
		Password: "12345678",
		Status:   database.UserStatusPending,
	})
	Expect(err).To(BeNil())

	// verify user
	path := "/auth/verify?code=" + code
	resp := execRequest(users[0].Id, http.MethodPost, path, "")
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	_, err = io.ReadAll(resp.Body)
	// check user status is active
	user, err := db.GetUser("user2@example.com")
	Expect(err).To(BeNil())
	Expect(user.Email).To(Equal("user2@example.com"))
	Expect(user.Status).To(Equal(database.UserStatusActive))
	err = resp.Body.Close()
	Expect(err).To(BeNil())
}

func TestRecover(t *testing.T) {
	initialize(t)

	path := "/auth/recover"
	body := fmt.Sprintf(`{"email": "%s", "recaptcha_token": "123456"}`, users[0].Email)
	resp := execRequest("", http.MethodPost, path, body)
	b, err := io.ReadAll(resp.Body)
	Expect(err).To(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK), string(b))
	err = resp.Body.Close()
	Expect(err).To(BeNil())

	// should have a verification of type recover
	_, err = db.GetVerification(users[0].Id, database.VerificationTypeRecover)
	Expect(err).To(BeNil())

	// duplicate request
	resp = execRequest("", http.MethodPost, path, body)
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	_, err = io.ReadAll(resp.Body)
	Expect(err).To(BeNil())
	err = resp.Body.Close()
	Expect(err).To(BeNil())

	// should overwrite existing code
	_, err = db.GetVerification(users[0].Id, database.VerificationTypeRecover)
	Expect(err).To(BeNil())
}

func TestReset(t *testing.T) {
	initialize(t)

	path := "/auth/recover"
	body := fmt.Sprintf(`{"email": "%s", "recaptcha_token": "123456"}`, users[0].Email)
	resp := execRequest("", http.MethodPost, path, body)
	b, err := io.ReadAll(resp.Body)
	Expect(err).To(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK), string(b))
	err = resp.Body.Close()
	Expect(err).To(BeNil())

	code, err := db.GetVerification(users[0].Id, database.VerificationTypeRecover)
	Expect(err).To(BeNil())

	path = "/auth/reset"
	body = fmt.Sprintf(`{"password": "password2", "code": "%s", "recaptcha_token": "123456"}`, code)
	resp = execRequest(users[0].Id, http.MethodPatch, path, body)
	Expect(resp.StatusCode).To(Equal(http.StatusAccepted))
	_, err = io.ReadAll(resp.Body)
	Expect(err).To(BeNil())
	err = resp.Body.Close()
	Expect(err).To(BeNil())

	_, err = login(users[0].Email, "password2")
	Expect(err).To(BeNil())
}

func TestMain(m *testing.M) {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = true
	client = &http.Client{Transport: t, Timeout: time.Minute}
	recaptchaServer.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		resp := recaptcha.Response{
			Success:     true,
			Score:       1.0,
			Action:      "login",
			ChallengeTS: time.Now(),
			Hostname:    "localhost:8080",
			ErrorCodes:  nil,
		}
		b, _ := jsoniter.Marshal(&resp)
		if _, err := writer.Write(b); err != nil {
			panic(err)
		}

	}
	go recaptchaServer.Start()
	var err error
	db, err = database.Connect(&database.Config{connectionStr})
	if err != nil {
		panic(err)
	}
	s := NewServer(
		&serverConfig,
		db,
		&mailer.Mock{
			SendEMailVerificationFunc: func(toName string, toEmail string, code string) error {
				return nil
			},
			SendPasswordResetFunc: func(toName string, toEmail string, code string) error {
				return nil
			},
		},
		zap.L(),
	)
	go func() {
		err := s.ListenAndServer()
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	err = db.Clear(true)
	if err != nil {
		panic(err)
	}
	tokens = make(map[database.ObjectId]string)
	for i := range users {
		users[i].Id, _, err = db.AddUser(database.NewUser{Email: users[i].Email, Password: users[i].Email, Status: database.UserStatusActive})
		if err != nil {
			panic(err)
		}
		token, err := login(users[i].Email, users[i].Password)
		if err != nil {
			panic(err)
		}
		tokens[users[i].Id] = token
	}
	m.Run()
	err = s.Shutdown()
	if err != nil {
		panic(err)
	}
	err = db.Close()
	if err != nil {
		panic(err)
	}
}

func generateURL(path string) string {
	return "http://" + serverConfig.BindAddress + path
}

func login(user string, password string) (string, error) {
	url := generateURL("/auth/login")
	body := strings.NewReader(fmt.Sprintf(`{"email":"%s", "password": "%s", "recaptcha_token": "123456"}`, user, password))
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return "", err
	}
	req.Close = true
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("3")
		return "", err
	}

	loginResp := make(map[string]interface{})
	err = json.Unmarshal(respBody, &loginResp)
	if err != nil {
		return "", err
	}
	if loginResp["code"].(float64) != 200 {
		fmt.Println(loginResp)
		return "", errors.New("login failed")
	}
	return loginResp["token"].(string), nil
}

func initialize(t *testing.T) {
	RegisterTestingT(t)
	err := db.Clear(false)
	Expect(err).To(BeNil())
}

func execRequest(userId database.ObjectId, method string, path string, body string) *http.Response {
	url := generateURL(path)
	reqBody := strings.NewReader(body)
	req, err := http.NewRequest(method, url, reqBody)
	Expect(err).To(BeNil())
	req.Close = true
	req.Header.Add("Content-Type", "application/json")
	if userId != "" {
		req.Header.Add("Authorization", "Bearer "+tokens[userId])
	}
	resp, err := client.Do(req)
	Expect(err).To(BeNil())
	return resp
}
