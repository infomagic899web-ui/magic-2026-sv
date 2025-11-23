package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Image struct {
	URL string `json:"url" bson:"url"`
}

type Artist struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

type ExternalUrl struct {
	Spotify string `json:"spotify,omitempty" bson:"spotify,omitempty"`
}

type Album struct {
	AlbumType            string      `json:"album_type" bson:"album_type"`
	Href                 string      `json:"href" bson:"href"`
	ID                   string      `json:"id" bson:"id"`
	Artists              []Artist    `json:"artists" bson:"artists"`
	ExternalUrls         ExternalUrl `json:"external_urls" bson:"external_urls"`
	Images               []Image     `json:"images" bson:"images"`
	IsPlayable           bool        `json:"is_playable" bson:"is_playable"`
	Name                 string      `json:"name" bson:"name"`
	ReleaseDate          string      `json:"release_date" bson:"release_date"`
	ReleaseDatePrecision string      `json:"release_date_precision" bson:"release_date_precision"`
	TotalTracks          int         `json:"total_tracks" bson:"total_tracks"`
	Type                 string      `json:"type" bson:"type"`
	URI                  string      `json:"uri" bson:"uri"`
}

type Track struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Album        Album              `json:"album" bson:"album"`
	Artists      []Artist           `json:"artists" bson:"artists"`
	RequestCount int                `json:"requestCount" bson:"requestCount"`
	RequestedBy  []string           `json:"requestedBy" bson:"requestedBy"`
	ExternalUrls ExternalUrl        `json:"external_urls" bson:"external_urls"`
	IsLocal      bool               `json:"is_local" bson:"is_local"`
	TrackNumber  int                `json:"track_number" bson:"track_number"`
	CanVote      bool               `json:"can_vote" bson:"can_vote"` // <--- New field
	Created_at   primitive.DateTime `json:"created_at" bson:"created_at"`
	Updated_at   primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

const RequestLimitPer72Hours = 3
