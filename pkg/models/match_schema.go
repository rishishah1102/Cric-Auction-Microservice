package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Match struct {
	Id                primitive.ObjectID `bson:"_id" json:"_id"`
	Matches           []int              `bson:"matches" json:"matches"`
	PrevX1            bool               `bson:"prevX1" json:"prevX1"`
	CurrentX1         bool               `bson:"currentX1" json:"currentX1"`
	NextX1            bool               `bson:"nextX1" json:"nextX1"`
	EarnedPoints      int                `bson:"earnedPoints" json:"earnedPoints"`
	BenchedPoints     int                `bson:"benchedPoints" json:"benchedPoints"`
	TotalPoints       int                `bson:"totalPoints" json:"totalPoints"`
	PrevTotalPoints   int                `bson:"prevTotalPoints" json:"prevTotalPoints"`
	PrevEarnedPoints  int                `bson:"prevEarnedPoints" json:"prevEarnedPoints"`
	PrevBenchedPoints int                `bson:"prevBenchedPoints" json:"prevBenchedPoints"`
}
