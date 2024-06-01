package database

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"time"
)

func randomString(n int) string {
	const (
		letterBytes   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

var (
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrNotFound       = errors.New("not found")
	ErrInvalid        = errors.New("invalid operation")
	ErrUnauthorized   = errors.New("authorization failed")
)

type Database struct {
	db *sql.DB
}

func Connect(config *Config) (*Database, error) {
	db, err := sql.Open("mysql", config.ConnectionString)
	if err != nil {
		return nil, parseError(err)
	}
	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}

func (db *Database) Clear(removeUsers bool) error {
	var err error
	err = db.withTransaction(func(t *sql.Tx) error {
		if err := deleteVerifications(t); err != nil {
			return err
		}
		if removeUsers {
			if err := deleteUsers(t); err != nil {
				return err
			}
		}
		return nil
	})
	return parseError(err)
}

func (db *Database) AddUser(u NewUser) (ObjectId, string, error) {
	hash, err := HashPassword(u.Password)
	if err != nil {
		return EmptyObjectId, "", err
	}
	u.Password = hash
	userId := NewObjectId()
	code := randomString(50)
	err = db.withTransaction(func(t *sql.Tx) error {
		if err := addUser(t, userId, u); err != nil {
			return err
		}
		if u.Status == UserStatusPending {
			err := setVerification(t, userId, Verification{Code: code, Type: VerificationTypeSignup})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return userId, code, parseError(err)
}

func (db *Database) Verify(code string) error {
	err := db.applyVerifiedAction(code, VerificationTypeSignup, func(t *sql.Tx, userId ObjectId) error {
		return setUserStatus(t, userId, UserStatusActive)
	})
	return parseError(err)
}

func (db *Database) SetRecoveryCode(userId ObjectId) (string, error) {
	code := randomString(50)
	err := db.withTransaction(func(t *sql.Tx) error {
		err := setVerification(t, userId, Verification{Type: VerificationTypeRecover, Code: code})
		return err
	})
	if err != nil {
		return "", parseError(err)
	}
	return code, nil
}

func (db *Database) ResetPassword(code string, newPassword string) error {
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	err = db.applyVerifiedAction(code, VerificationTypeRecover, func(t *sql.Tx, userId ObjectId) error {
		return setUserPassword(t, userId, hash)
	})
	return parseError(err)
}

func (db *Database) GetUser(name string) (User, error) {
	u, err := db.getUser(name)
	return u, parseError(err)
}

func (db *Database) DeleteUser(name string) error {
	return parseError(db.deleteUser(name))
}

func (db *Database) GetVerification(userId ObjectId, verificationType VerificationType) (string, error) {
	code, err := db.getVerification(userId, verificationType)
	if err != nil {
		return "", parseError(err)
	}
	return code, nil
}
