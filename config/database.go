package config

import (
	"gopkg.in/mgo.v2"
)

// NewSession return the connection with mongoDB
func NewSession() (*mgo.Session, error) {
	s, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}
	return s, nil
}
