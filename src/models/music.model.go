package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Music struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Music_image    string             `json:"music_image"`
	Title          string             `json:"title"`
	Artist         []string           `json:"artist"`
	Album          string             `json:"album"`
	Votes          int                `json:"votes"`
	Upcoming_votes int                `json:"upcoming_votes"`
	Music_type     string             `json:"music_type"`
	Music_url      string             `json:"music_url"`
	PeekDate       primitive.DateTime `json:"peek_date"`
	Weeks          int                `json:"weeks"`
	Created_at     primitive.DateTime `json:"created_at"`
	Updated_at     primitive.DateTime `json:"updated_at"`
}


// EVERY monday listeners can vote but it will not adding it to the actual vote instead, it will throw to upcoming_votes and after friday the vote is close therefore the votes from the upcoming_votes will finally add to the actual vote