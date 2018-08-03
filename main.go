package main

import (
	"log"
	"net/http"

	"github.com/bookfire/config"
	"github.com/bookfire/controller"
	"github.com/gorilla/mux"
)

func main() {
	conn, err := config.NewSession()

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	router := mux.NewRouter()
	artistHandler := controller.NewArtistHandler(conn)

	router.HandleFunc("/artists/", artistHandler.Create).Methods("POST")
	router.HandleFunc("/artists/", artistHandler.Read).Methods("GET")
	router.HandleFunc("/artists/{id}", artistHandler.FindByID).Methods("GET")
	router.HandleFunc("/artists/{id}", artistHandler.Update).Methods("PUT")
	router.HandleFunc("/artists/{id}", artistHandler.Delete).Methods("DELETE")

	cors := config.SetupCors()
	log.Fatal(http.ListenAndServe(":8081", cors(router)))
}
