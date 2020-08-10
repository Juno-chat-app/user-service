package entity

import "time"

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
)

type User struct {
	UserName    string        `json:"user-name"`
	Password    string        `json:"password"`
	Status      *UserStatus   `json:"status"`
	ContactInfo *ContactInfo  `json:"contact-info"`
	Permissions []*Permission `json:"permissions"`
}

type UserStatus struct {
	Status         Status     `json:"user-status"`
	ActivationCode string     `json:"activation-code"`
	UpdatedAt      *time.Time `json:"updated-at"`
}

type ContactInfo struct {
	Mobile string `json:"mobile"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
}
