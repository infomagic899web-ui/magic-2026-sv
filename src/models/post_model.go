package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title        string             `json:"title" validate:"required,min=2,max=100"`
	Details      []string           `json:"details" validate:"required,min=2,max=100"`
	Image        string             `json:"image"`
	Link         string             `json:"link"`
	Social_media string             `json:"social_media" validate:"required,oneof=tiktok instagram facebook"`
	Created_at   primitive.DateTime `json:"created_at"`
	Updated_at   primitive.DateTime `json:"updated_at"`
}
