package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	ReleaseDate time.Time  `json:"release_date"`
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
		ReleaseDate: playerResp.Item.Album.ReleaseDateTime(),
	}
}

func playerHandler(w http.ResponseWriter, r *http.Request) {
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

type playInfoRequest struct {
	URI spotify.URI `json:"uri"`
}

func playMusicHandler(w http.ResponseWriter, r *http.Request) {
	var playInfo playInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&playInfo); err != nil {
		log.Println(err)
		return
	}
	bearer := r.Header.Get("Authorization")
	token := new(oauth2.Token)
	token.AccessToken = bearer
	client := spotify.Authenticator{}.NewClient(token)

	playOptions := spotify.PlayOptions{}
	if len(playInfo.URI) > 0 {
		playOptions.PlaybackContext = &playInfo.URI
	}
	if err := client.PlayOpt(&playOptions); err != nil {
		log.Println(err)
		return
	}
}

func pauseMusicHandler(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := new(oauth2.Token)
	token.AccessToken = bearer
	client := spotify.Authenticator{}.NewClient(token)

	if err := client.Pause(); err != nil {
		log.Println(err)
		return
	}
}

func nextMusicHandler(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := new(oauth2.Token)
	token.AccessToken = bearer
	client := spotify.Authenticator{}.NewClient(token)

	if err := client.Next(); err != nil {
		log.Println(err)
		return
	}
}

func prevMusicHandler(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := new(oauth2.Token)
	token.AccessToken = bearer
	client := spotify.Authenticator{}.NewClient(token)

	if err := client.Previous(); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/player", playerHandler).Methods("GET")
	r.HandleFunc("/player/play", playMusicHandler).Methods("POST")
	r.HandleFunc("/player/pause", pauseMusicHandler).Methods("POST")
	r.HandleFunc("/player/next", nextMusicHandler).Methods("POST")
	r.HandleFunc("/player/prev", prevMusicHandler).Methods("POST")

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(r)))
}
