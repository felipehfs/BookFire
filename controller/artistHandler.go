// Package controller represents
// the handler which operates
// the data on the routes
package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bookfire/model"
	"gopkg.in/mgo.v2"
)

// ArtistHandler represents the controller of artist
type ArtistHandler struct {
	db *mgo.Session
}

// NewArtistHandler returns the a new instance
func NewArtistHandler(d *mgo.Session) *ArtistHandler {
	return &ArtistHandler{d}
}

// Create - adds new resource to database
func (handle ArtistHandler) Create(w http.ResponseWriter, req *http.Request) {
	var artist model.Artist
	if err := json.NewDecoder(req.Body).Decode(&artist); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
	dao := model.NewArtistDAO(handle.db)

	if err := dao.Create(artist); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("%s - %s - %s\n", req.Method, req.Host, req.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

// Read retrives all data to JSON
func (handle ArtistHandler) Read(w http.ResponseWriter, req *http.Request) {
	dao := model.NewArtistDAO(handle.db)
	results, err := dao.Read()

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("%s - %s - %s\n", req.Method, req.Host, req.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// Update represents the controller that makes
// the changes of a artist
func (handle ArtistHandler) Update(w http.ResponseWriter, req *http.Request) {
	var artist model.Artist
	if err := json.NewDecoder(req.Body).Decode(&artist); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	params := mux.Vars(req)

	dao := model.NewArtistDAO(handle.db)

	if err := dao.Update(params["id"], artist); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("%s - %s - %s\n", req.Method, req.Host, req.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

// FindByID searches the artist on the database
func (handle ArtistHandler) FindByID(w http.ResponseWriter, req *http.Request) {
	param := mux.Vars(req)

	artist, err := model.NewArtistDAO(handle.db).FindByID(param["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	log.Printf("%s - %s - %s\n", req.Method, req.Host, req.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

// Delete - removes the resource from the database
func (handle ArtistHandler) Delete(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	if err := model.NewArtistDAO(handle.db).Delete(params["id"]); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("%s - %s - %s\n", req.Method, req.Host, req.URL.Path)
	w.Write([]byte("success!"))
}
