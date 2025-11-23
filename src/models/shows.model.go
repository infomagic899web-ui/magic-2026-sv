package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Shows struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Show_name    *string            `json:"show_name" validate:"required,min=2,max=100"`
	Show_host    *string            `json:"show_host" validate:"required,min=2,max=100"`
	Show_desc    []string           `json:"show_desc" validate:"required,min=2,max=100"`
	Show_day     *string            `json:"show_day" validate:"required,min=2,max=100"`
	Show_time    *string            `json:"show_time" validate:"required,min=2,max=100"`
	Show_Image   string             `json:"show_image"`
	Social_media []string           `json:"social_media"`
	Is_active    bool               `json:"is_active"`
	Created_at   primitive.DateTime `json:"created_at"`
	Updated_at   primitive.DateTime `json:"updated_at"`
}
