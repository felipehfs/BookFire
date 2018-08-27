// Package model represents the business api logic
package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserDAO represents the dao of user
type UserDAO struct {
	conn *mgo.Session
}

// NewUserDAO generates the new instance of UserDAO
func NewUserDAO(c *mgo.Session) *UserDAO {
	return &UserDAO{c}
}

// getColletion retrieves the collection of dao to manage
func (u UserDAO) getCollection() *mgo.Collection {
	return u.conn.DB("bookfire").C("users")
}

// Create insert data into database
func (u UserDAO) Create(user User) error {
	return u.getCollection().Insert(&user)
}

// Find searches the data into database
func (u UserDAO) Find(login string, password string) (*User, error) {
	var result User
	err := u.getCollection().Find(bson.M{"login": login, "password": password}).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
