package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Match struct {
	Id                primitive.ObjectID `bson:"_id" json:"_id"`
	Match1            int                `bson:"match1" json:"match1"`
	Match2            int                `bson:"match2" json:"match2"`
	Match3            int                `bson:"match3" json:"match3"`
	Match4            int                `bson:"match4" json:"match4"`
	Match5            int                `bson:"match5" json:"match5"`
	Match6            int                `bson:"match6" json:"match6"`
	Match7            int                `bson:"match7" json:"match7"`
	Match8            int                `bson:"match8" json:"match8"`
	Match9            int                `bson:"match9" json:"match9"`
	Match10           int                `bson:"match10" json:"match10"`
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
