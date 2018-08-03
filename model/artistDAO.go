//Package model represents the
// bussiness logic of the Rest api
package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ArtistDAO represent the data access object
// of Artist and makes all operation with database
type ArtistDAO struct {
	conn *mgo.Session
}

// NewArtistDAO the contruct of ArtistDao
// returns the a new instance
func NewArtistDAO(c *mgo.Session) *ArtistDAO {
	return &ArtistDAO{c}
}

// The getColletion retrieves the collection of dao to manage
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

// Update change the data of artist by ID
func (dao ArtistDAO) Update(id string, a Artist) error {
	return dao.getCollection().Update(bson.M{"_id": bson.ObjectIdHex(id)}, a)
}

// FindByName search the artist by name
func (dao ArtistDAO) FindByName(name string) (*Artist, error) {
	var result Artist
	if err := dao.getCollection().Find(bson.M{"name": name}).One(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// FindByID search the artist by ID
func (dao ArtistDAO) FindByID(id string) (*Artist, error) {
	var searched Artist
	err := dao.getCollection().Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&searched)
	if err != nil {
		return nil, err
	}
	return &searched, nil
}

// Delete removes the id from artist
func (dao ArtistDAO) Delete(id string) error {
	return dao.getCollection().Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
