package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func playlistHandler(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := new(oauth2.Token)
	token.AccessToken = bearer
	client := spotify.Authenticator{}.NewClient(token)
	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(playlists)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/playlist", playlistHandler).Methods("GET")

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(r)))
}
