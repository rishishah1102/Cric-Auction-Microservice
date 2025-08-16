package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auction struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	AuctionImg  string             `bson:"auctionImg" json:"auctionImg"`
	AuctionName string             `bson:"auctionName" json:"auctionName"`
	CreatedBy   string             `bson:"createdBy" json:"createdBy"`
	JoinedBy    []struct {
		Email string `bson:"email" json:"email"`
		Name  string `bson:"name" json:"name"`
	} `bson:"joinedBy" json:"joinedBy"`
	PointsTableChecked bool      `bson:"pointsTableChecked" json:"pointsTableChecked"`
	CreatedAt          time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time `bson:"updatedAt" json:"updatedAt"`
}
