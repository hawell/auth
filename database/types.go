package database

import (
	"github.com/google/uuid"
)

type ObjectId string

const EmptyObjectId ObjectId = ""

func NewObjectId() ObjectId {
	return ObjectId(uuid.New().String())
}

type UserStatus string

type User struct {
	Id       ObjectId
	Email    string
	Password string
	Status   UserStatus
}

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
	UserStatusPending  UserStatus = "pending"
)

type NewUser struct {
	Email    string
	Password string
	Status   UserStatus
}

type VerificationType string

const (
	VerificationTypeSignup  VerificationType = "signup"
	VerificationTypeRecover VerificationType = "recover"
)

type Verification struct {
	Code string
	Type VerificationType
}
