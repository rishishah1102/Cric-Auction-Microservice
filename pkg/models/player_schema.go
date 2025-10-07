package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	Id                primitive.ObjectID `bson:"_id" json:"_id"`
	AuctionId         primitive.ObjectID `bson:"auction_id" json:"auction_id" binding:"required"`
	PlayerNumber      int                `bson:"player_number" json:"player_number" binding:"required"`
	PlayerName        string             `bson:"player_name" json:"player_name" binding:"required"`
	Country           string             `bson:"country,omitempty" json:"country,omitempty"`
	Role              string             `bson:"role" json:"role" binding:"required"`
	PrevTeam          string             `bson:"prev_team" json:"prev_team"`
	CurrentTeam       string             `bson:"current_team" json:"current_team"`
	Hammer            string             `bson:"hammer" json:"hammer"`
	BasePrice         float64            `bson:"base_price" json:"base_price" binding:"required"`
	SellingPrice      float64            `bson:"selling_price" json:"selling_price"`
	IPLTeam           string             `bson:"ipl_team,omitempty" json:"ipl_team,omitempty"`
	PrevFantasyPoints int                `bson:"prev_fantasy_points,omitempty" json:"prev_fantasy_points,omitempty"`
	Bids              []Bids             `bson:"bids" json:"bids"`
	Match             primitive.ObjectID `bson:"match,omitempty" json:"match,omitempty"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

type Bids struct {
	TeamName string  `bson:"team_name" json:"team_name"`
	Bid      float64 `bson:"bid" json:"bid"`
}
