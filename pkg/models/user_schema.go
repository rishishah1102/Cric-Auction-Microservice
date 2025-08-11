package models

import (
	"time"

	"github.com/google/uuid"
)

// User is the struct for auction users
type User struct {
	UUID      uuid.UUID `db:"uuid" json:"uuid"`
	Id        int64     `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Mobile    string    `db:"mobile" json:"mobile"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
