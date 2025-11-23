package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movies struct {
	ID                      primitive.ObjectID `bson:"_id" json:"id"`
	Title                   string             `json:"title" validate:"required,min=2,max=100"`
	Content                 []string           `json:"content" validate:"required,min=2,max=100"`
	Movie_Image             string             `json:"movie_image"`
	Category                []string           `json:"category" validate:"required,min=2,max=100"`
	Social_media            []string           `json:"social_media"`
	Sponsors                []string           `json:"sponsors"`
	Cast                    []string           `json:"cast"`
	Movie_file              string             `json:"movie_file"`
	Release_date            string             `json:"release_date" validate:"required,min=2,max=100"`
	Directed_by             string             `json:"directed_by" validate:"required,min=2,max=100"`
	Rated                   string             `json:"rated" validate:"required,min=2,max=100"`
	Location_cinema         string             `json:"location_cinema" validate:"required,min=2,max=100"`
	Advanced_screening_date string             `json:"advanced_screening_date" validate:"required,min=2,max=100"`
	Screening_time          string             `json:"screening_time" validate:"required,min=2,max=100"`
	Rating                  float64            `json:"rating"`
	Created_at              primitive.DateTime `json:"created_at"`
	Updated_at              primitive.DateTime `json:"updated_at"`
}
