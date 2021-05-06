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

type player struct {
	IsPlaying   bool       `json:"is_playing"`
	AlbumName   string     `json:"album_name"`
	ArtistsName []string   `json:"artists_name"`
	MusicName   string     `json:"music_name"`
	ID          spotify.ID `json:"ID"`
}

func reducePlayer(playerResp *spotify.CurrentlyPlaying) player {
	var artists []string

	for _, p := range playerResp.Item.Artists {
		artists = append(artists, p.Name)
	}

	return player{
		IsPlaying:   playerResp.Playing,
		AlbumName:   playerResp.Item.Album.Name,
		ArtistsName: artists,
		MusicName:   playerResp.Item.Name,
		ID:          playerResp.Item.ID,
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := new(oauth2.Token)
	token.AccessToken = bearer
	client := spotify.Authenticator{}.NewClient(token)
	player, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		log.Println(err)
		return
	}

	reducedPlayer := reducePlayer(player)

	json.NewEncoder(w).Encode(reducedPlayer)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/player", userHandler).Methods("GET")

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(r)))
}
