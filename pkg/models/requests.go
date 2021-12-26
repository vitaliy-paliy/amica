package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type FriendRequest struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Sender    *User              `json:"sender"`
	Recipient *User              `json:"recipient"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp,omitempty"`
}

func (req *FriendRequest) Invert() *FriendRequest {
	return &FriendRequest{Sender: req.Recipient, Recipient: req.Sender}
}

func (req *FriendRequest) SelfRequest() bool {
	return req.Sender.ID == req.Recipient.ID
}

func (req *FriendRequest) ToSlice() *[]*User {
	return &[]*User{req.Sender, req.Recipient}
}

func (req *FriendRequest) ToUserDoc() *primitive.D {
	return &primitive.D{{"sender", req.Sender}, {"recipient", req.Recipient}}
}
