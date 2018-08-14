// Package config make the setup of all application
package config

import (
	"gopkg.in/mgo.v2"
)

// NewSession returns the connection with mongoDB
func NewSession() (*mgo.Session, error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}
	return session, nil
}
