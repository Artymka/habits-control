package models

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type UserCreate struct {
	Email        string
	PasswordHash string
}
