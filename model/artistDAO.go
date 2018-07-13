package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ArtistDAO
type ArtistDAO struct {
	conn *mgo.Session
}

// NewArtistDAO returns the a new instance
func NewArtistDAO(c *mgo.Session) *ArtistDAO {
	return &ArtistDAO{c}
}

func (dao ArtistDAO) getCollection() *mgo.Collection {
	return dao.conn.DB("bookfire").C("artists")
}

// Create insert data into database
func (dao ArtistDAO) Create(a Artist) error {
	return dao.getCollection().Insert(a)
}

// Read retrieves all data
func (dao ArtistDAO) Read() ([]Artist, error) {
	var results []Artist
	if err := dao.getCollection().Find(bson.M{}).All(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// Update change the data of artist
func (dao ArtistDAO) Update(id string, a Artist) error {
	return dao.getCollection().Update(bson.M{"_id": bson.ObjectIdHex(id)}, a)
}

// Delete removes the id from artist
func (dao ArtistDAO) Delete(id string) error {
	return dao.getCollection().Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
