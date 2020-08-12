package entity

import (
	"time"
)

type (
	Status string
	Path   string
)

const (
	Active          Status = "active"
	Inactive        Status = "inactive"
	DocumentVersion string = "v0.0.1"

	// The paths to access data in database
	UserNamePath  Path = "user-name"
	EmailPath     Path = "contact-info.email"
	StatusPath    Path = "status.user-status"
	PasswordPath  Path = "password"
	UserIdPath    Path = "user-id"
	DeletedAtPath Path = "deleted-at"
)

type User struct {
	UserName        string        `bson:"user-name"`
	Password        string        `bson:"password"`
	UserId          string        `bson:"user-id"`
	Status          *UserStatus   `bson:"status"`
	ContactInfo     *ContactInfo  `bson:"contact-info"`
	Permissions     []*Permission `bson:"permissions"`
	CreatedAt       *time.Time    `bson:"created-at"`
	UpdatedAt       *time.Time    `bson:"updated-at"`
	DeletedAt       *time.Time    `bson:"deleted-at"`
	DocumentVersion string        `bson:"document-version"`
}

type UserStatus struct {
	Status         Status     `bson:"user-status"`
	ActivationCode string     `bson:"activation-code"`
	UpdatedAt      *time.Time `bson:"updated-at"`
}

type ContactInfo struct {
	Mobile string `bson:"mobile"`
	Phone  string `bson:"phone"`
	Email  string `bson:"email"`
}
