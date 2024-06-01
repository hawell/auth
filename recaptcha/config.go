package recaptcha

type Config struct {
	Bypass    bool   `json:"bypass"`
	SecretKey string `json:"secret_key"`
	Server    string `json:"server"`
}
