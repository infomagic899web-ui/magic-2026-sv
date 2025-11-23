package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MagicVideos struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Title      string             `bson:"title" json:"title"`
	Video_url  string             `bson:"video_url" json:"video_url"`
	Desc       []string           `bson:"desc" json:"desc"`
	Show_name  string             `bson:"show_name" json:"show_name"`
	Provider   string             `bson:"provider" json:"provider"`
	Thumbnail  string             `bson:"thumbnail" json:"thumbnail"`
	Socials    []*string          `bson:"socials" json:"socials"`
	Date       string             `bson:"date" json:"date"`
	Created_at primitive.DateTime `bson:"created_at" json:"created_at"`
	Updated_at primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
