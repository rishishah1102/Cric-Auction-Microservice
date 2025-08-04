package models

import (
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
	CreatedAt          int64 `bson:"createdAt" json:"createdAt"`
	PointsTableChecked bool  `bson:"pointsTableChecked" json:"pointsTableChecked"`
}
