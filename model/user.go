package model

import (
	"gopkg.in/mgo.v2/bson"
)

// User represents the login entity
type User struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Login    string        `json:"login"`
	Password string        `json:"password"`
}
