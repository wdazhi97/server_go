package entity

import "time"

type Friendship struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	FriendID  string    `bson:"friend_id" json:"friend_id"`
	Status    string    `bson:"status" json:"status"` // pending, accepted, blocked
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}