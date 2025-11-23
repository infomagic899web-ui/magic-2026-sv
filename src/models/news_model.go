package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Title           string             `json:"title" validate:"required,min=2,max=100"`
	NormalizedTitle string             `bson:"normalized_title" json:"normalized_title"`
	Content         []string           `json:"content" validate:"required,min=2,max=100"`
	News_Image      string             `json:"news_image"`
	Category        string             `json:"category" validate:"required,min=2,max=100,eq=News|eq=Sports|eq=celebrities|eq=Music|eq=Movies"`
	Social_media    []string           `json:"social_media"`
	Writer          string             `json:"writer" validate:"required,min=2,max=100"`
	Status          string             `json:"status" validate:"required,min=2,max=100,eq=pending|eq=approved|eq=rejected"`
	Sources         []string           `json:"sources" validate:"required,min=2,max=100"`
	DateTime        string             `json:"datetime" validate:"required,min=2,max=100"`
	ImageReference  string             `json:"imagereference"`
	Created_at      primitive.DateTime `json:"created_at"`
	Updated_at      primitive.DateTime `json:"updated_at"`
}
