package entity

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id,omitempty"`
	Username  string    `bson:"username"`
	Password  string    `bson:"password"`
	Email     string    `bson:"email"`
	Online    bool      `bson:"online"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	LastSeen  time.Time `bson:"last_seen"`
}