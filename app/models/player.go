package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Credits int                `bson:"credits"`
	Status  string             `bson:"status"`
}

type GameOutcome struct {
	PlayerID  primitive.ObjectID `bson:"player_id"`
	Timestamp time.Time          `bson:"timestamp"`
	Result    []string           `bson:"result"`
	WinAmount int                `bson:"win_amount"`
	TotalBet  int                `bson:"total_bet"`
}

type RTPStatistic struct {
	ID          string    `bson:"_id,omitempty"`
	TotalPlays  int       `bson:"total_plays"`
	TotalWins   int       `bson:"total_wins"`
	TotalLosses int       `bson:"total_losses"`
	TotalPayout int       `bson:"total_payout"`
	LastUpdated time.Time `bson:"last_updated"`
}
