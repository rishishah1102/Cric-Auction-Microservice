package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auction struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	AuctionName  string             `bson:"auction_name" json:"auction_name"`
	AuctionImage string             `bson:"auction_image" json:"auction_image"`
	CreatedBy    string             `bson:"created_by" json:"created_by"`
	AuctionDate  time.Time          `bson:"auction_date" json:"auction_date"`
	IsIPLAuction bool               `bson:"is_ipl_auction" json:"is_ipl_auction"`
	JoinedBy     []string           `bson:"joined_by" json:"joined_by"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
