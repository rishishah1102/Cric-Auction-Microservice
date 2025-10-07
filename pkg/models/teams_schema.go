package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID         primitive.ObjectID   `bson:"_id" json:"id"`
	TeamName   string               `bson:"team_name" json:"team_name"`
	TeamImage  string               `bson:"team_image" json:"team_image"`
	AuctionId  primitive.ObjectID   `bson:"auction_id" json:"auction_id"`
	TeamOwners []string             `bson:"team_owners" json:"team_owners"`
	Squad      []primitive.ObjectID `bson:"squad" json:"squad"`
	CreatedAt  time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time            `bson:"updated_at" json:"updated_at"`
}
