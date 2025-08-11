package models

import (
	"time"

	"github.com/google/uuid"
)

// Profile is the struct for profile of auction users
type Profile struct {
	UUID         uuid.UUID `db:"uuid" json:"uuid"`
	Id           int64     `db:"id" json:"id"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	ImageURL     string    `db:"image_url" json:"image_url"`
	Role         string    `db:"role" json:"role"`
	BattingHand  string    `db:"batting_hand" json:"batting_hand"`
	BattingOrder string    `db:"batting_order" json:"batting_order"`
	BattingStyle string    `db:"batting_style" json:"batting_style"`
	BowlingArm   string    `db:"bowling_arm" json:"bowling_arm"`
	BowlingType  string    `db:"bowling_type" json:"bowling_type"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
