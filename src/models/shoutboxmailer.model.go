package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// RequestShoutbox is the model for storing shoutbox requests.
type RequestShoutbox struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" validate:"required"`
	School_name    string             `json:"school_name" validate:"required"`
	Email          string             `json:"email" validate:"required,email"`
	Position       string             `json:"position" validate:"required"`
	Contact        string             `json:"contact" validate:"required,max=11"`
	School_contact string             `json:"school_contact" validate:"required,max=11"`
	Organization   string             `json:"organization" validate:"required"`
	Title          string             `json:"title" validate:"required"`
	Event_date     string             `json:"event_date" validate:"required"`
	Radio_spiel    string             `json:"radio_spiel" validate:"required,min=100,max=300"`
	Created_at     primitive.DateTime `json:"created_at"`
	Updated_at     primitive.DateTime `json:"updated_at"`
}
