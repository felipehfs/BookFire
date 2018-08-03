// Package model represents the bussiness logic
package model

import (
	"gopkg.in/mgo.v2/bson"
)

// Artist represent a person
type Artist struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string        `json:"name,omitmepty" bson:"name,omitempty" binding:"required"`
	Email       string        `json:"email,omitempty" bson:"email,omitempty" binding:"required"`
	Description string        `json:"description,omitempty" bson:"description,omitempty" binding:"required"`
}
