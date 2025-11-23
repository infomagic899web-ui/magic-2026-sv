package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VoteRecord struct {
	MusicID    primitive.ObjectID `bson:"music_id"`
	Date       string             `bson:"date"`
	IPAddress  string             `bson:"ip_address"`
	Votes      int                `bson:"votes"`
	LastVoteAt time.Time          `bson:"last_vote_at"` // New field
	CanVote    bool               `bson:"can_vote"`     // New field to track voting eligibility
	CanRequest bool               `bson:"can_request"`  // New field to track request eligibility
}
