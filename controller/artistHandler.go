package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bookfire/model"
	"gopkg.in/mgo.v2"
)

// ArtistHandler represent the controller of artist
type ArtistHandler struct {
	db *mgo.Session
}

// NewArtistHandler returns the a new instance
func NewArtistHandler(d *mgo.Session) *ArtistHandler {
	return &ArtistHandler{d}
}

// Create - add new resource to database
func (handle ArtistHandler) Create(w http.ResponseWriter, req *http.Request) {
	var artist model.Artist
	if err := json.NewDecoder(req.Body).Decode(&artist); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	dao := model.NewArtistDAO(handle.db)
	if err := dao.Create(artist); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

// Read retrives all data to JSON
func (handle ArtistHandler) Read(w http.ResponseWriter, req *http.Request) {
	dao := model.NewArtistDAO(handle.db)
	results, err := dao.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (handle ArtistHandler) Update(w http.ResponseWriter, req *http.Request) {
	var artist model.Artist
	if err := json.NewDecoder(req.Body).Decode(&artist); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	params := mux.Vars(req)
	err := model.NewArtistDAO(handle.db).Update(params["id"], artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

// Delete - removes the resource from the database
func (handle ArtistHandler) Delete(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if err := model.NewArtistDAO(handle.db).Delete(params["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte("success!"))
}
