package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	TeamId        primitive.ObjectID `bson:"_id" json:"_id"`
	AuctionId     primitive.ObjectID `bson:"auctionId" json:"auctionId"`
	TeamName      string             `bson:"teamName" json:"teamName"`
	TeamImg       string             `bson:"teamImg" json:"teamImg"`
	PartnerShip   []string           `bson:"partnerShip" json:"partnerShip"`
}
