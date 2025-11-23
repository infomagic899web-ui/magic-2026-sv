package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username"`
	Password    string             `bson:"password,omitempty" json:"-"`
	Email       string             `bson:"email" json:"email"`
	Role        string             `bson:"role" json:"role"`
	Avatar      string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	RSAPublic   string             `bson:"rsa_public,omitempty" json:"-"`
	RSAPrivate  string             `bson:"rsa_private,omitempty" json:"-"`
	IsVerified  bool               `bson:"is_verified" json:"is_verified"`
	SessionID   string             `bson:"session_id,omitempty" json:"-"`
	LastLoginAt time.Time          `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
