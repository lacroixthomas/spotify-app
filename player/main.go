package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type key int

var CLIENT_CONTEXT = key(1)

// spotifyClient interface of spotify client
type spotifyClient interface {
	PlayerCurrentlyPlaying() (*spotify.CurrentlyPlaying, error)
	PlayOpt(opt *spotify.PlayOptions) error
	Pause() error
	Next() error
	Previous() error
}

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
	client := r.Context().Value(CLIENT_CONTEXT).(spotifyClient)
	player, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		log.WithError(err).Error("playerHandler: could not get player currently playing")
		w.WriteHeader(http.StatusInternalServerError)
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
		log.WithError(err).Error("playMusicHandler: could not get decode play music request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := r.Context().Value(CLIENT_CONTEXT).(spotifyClient)
	playOptions := spotify.PlayOptions{}
	if len(playInfo.URI) > 0 {
		playOptions.PlaybackContext = &playInfo.URI
	}
	if err := client.PlayOpt(&playOptions); err != nil {
		log.WithError(err).Error("playMusicHandler: could not get play music")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func pauseMusicHandler(w http.ResponseWriter, r *http.Request) {
	client := r.Context().Value(CLIENT_CONTEXT).(spotifyClient)
	if err := client.Pause(); err != nil {
		log.WithError(err).Error("pauseMusicHandler: could not get pause music")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func nextMusicHandler(w http.ResponseWriter, r *http.Request) {
	client := r.Context().Value(CLIENT_CONTEXT).(spotifyClient)
	if err := client.Next(); err != nil {
		log.WithError(err).Error("nextMusicHandler: could not get play next music")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func prevMusicHandler(w http.ResponseWriter, r *http.Request) {
	client := r.Context().Value(CLIENT_CONTEXT).(spotifyClient)
	if err := client.Previous(); err != nil {
		log.WithError(err).Error("prevMusicHandler: could not get play previous music")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func tokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		token := new(oauth2.Token)
		token.AccessToken = bearer
		client := spotify.Authenticator{}.NewClient(token)
		ctx := r.Context()
		ctx = context.WithValue(ctx, CLIENT_CONTEXT, &client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

	contextedMux := tokenMiddleware(r)
	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(contextedMux)))
}
