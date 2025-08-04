package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Email        string             `bson:"email" json:"email"`
	Mobile       string             `bson:"mobile" json:"mobile"`
	ImgUrl       string             `bson:"imgUrl" json:"imgUrl"`
	// FirstName    string             `bson:"firstName" json:"firstName"`
	// LastName     string             `bson:"lastName" json:"lastName"`
	// Role         string             `bson:"role" json:"role"`
	// BowlingArm   string             `bson:"bowlingArm" json:"bowlingArm"`
	// BowlingType  string             `bson:"bowlingType" json:"bowlingType"`
	// BattingHand  string             `bson:"battingHand" json:"battingHand"`
	// BattingOrder string             `bson:"battingOrder" json:"battingOrder"`
	// BattingStyle string             `bson:"battingStyle" json:"battingStyle"`
}
