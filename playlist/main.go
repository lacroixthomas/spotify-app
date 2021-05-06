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

type playlistItem struct {
	Image     string      `json:"image"`
	Name      string      `json:"name"`
	OwnerName string      `json:"owner_name"`
	ID        spotify.ID  `json:"ID"`
	URI       spotify.URI `json:"uri"`
}

func reducePlaylist(playlistResp *spotify.SimplePlaylistPage) []playlistItem {
	var playlist []playlistItem

	for _, item := range playlistResp.Playlists {
		var image string
		if len(item.Images) > 0 {
			image = item.Images[0].URL
		}
		playlist = append(playlist, playlistItem{
			Image:     image,
			Name:      item.Name,
			OwnerName: item.Owner.DisplayName,
			ID:        item.ID,
			URI:       item.URI,
		})
	}
	return playlist

}

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

	reducedPlaylist := reducePlaylist(playlists)

	json.NewEncoder(w).Encode(reducedPlaylist)
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
