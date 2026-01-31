package entity

type Player struct {
	ID       string `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Rating   int32  `bson:"rating"`
	Status   string `bson:"status"` // waiting, matched, playing
}