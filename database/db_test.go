package database

import (
	. "github.com/onsi/gomega"
	"testing"
)

var (
	connectionStr = "admin:admin@tcp(127.0.0.1:3306)/auth"
	db            *Database
)

func TestConnect(t *testing.T) {
	RegisterTestingT(t)
	db, err := Connect(&Config{connectionStr})
	Expect(err).To(BeNil())
	err = db.Close()
	Expect(err).To(BeNil())
}

func TestUser(t *testing.T) {
	RegisterTestingT(t)
	err := db.Clear(true)
	Expect(err).To(BeNil())

	// add
	_, _, err = db.AddUser(NewUser{Email: "dbUser1", Password: "12345678", Status: UserStatusActive})
	Expect(err).To(BeNil())

	// get
	u, err := db.GetUser("dbUser1")
	Expect(err).To(BeNil())
	Expect(u.Email).To(Equal("dbUser1"))
	Expect(u.Status).To(Equal(UserStatusActive))

	// get non-existing user
	u, err = db.GetUser("dbUser2")
	Expect(err).To(Equal(ErrNotFound))

	// duplicate
	_, _, err = db.AddUser(NewUser{Email: "dbUser1", Password: "dbUser1", Status: UserStatusActive})
	Expect(err).To(Equal(ErrDuplicateEntry))

	// delete
	err = db.DeleteUser("dbUser1")
	Expect(err).To(BeNil())
	_, err = db.GetUser("dbUser1")
	Expect(err).NotTo(BeNil())

	// delete non-existing user
	err = db.DeleteUser("dbUser1")
	Expect(err).To(Equal(ErrNotFound))
}

func TestMain(m *testing.M) {
	db, _ = Connect(&Config{connectionStr})
	m.Run()
	_ = db.Close()
}
