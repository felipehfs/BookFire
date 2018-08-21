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
	userHandler := controller.NewUserHandler(conn)

	jwt := controller.EnabledJwt()
	adapt := config.ChainMiddleware(jwt, config.GzipHandler)

	router.HandleFunc("/artists", adapt(artistHandler.Create)).Methods("POST")
	router.HandleFunc("/artists", adapt(artistHandler.Read)).Methods("GET")
	router.HandleFunc("/artists/{id}", adapt(artistHandler.FindByID)).Methods("GET")
	router.HandleFunc("/artists/{id}", adapt(artistHandler.Update)).Methods("PUT")
	router.HandleFunc("/artists/{id}", adapt(artistHandler.Delete)).Methods("DELETE")
	router.HandleFunc("/register", userHandler.Create).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	cors := config.SetupCors()
	log.Fatal(http.ListenAndServe(":8081", cors(router)))
}
