package models

import "time"

type User struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// FirstName    string             `bson:"firstName" json:"firstName"`
	// LastName     string             `bson:"lastName" json:"lastName"`
	// Role         string             `bson:"role" json:"role"`
	// BowlingArm   string             `bson:"bowlingArm" json:"bowlingArm"`
	// BowlingType  string             `bson:"bowlingType" json:"bowlingType"`
	// BattingHand  string             `bson:"battingHand" json:"battingHand"`
	// BattingOrder string             `bson:"battingOrder" json:"battingOrder"`
	// BattingStyle string             `bson:"battingStyle" json:"battingStyle"`
}
