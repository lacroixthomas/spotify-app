package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// Key type of spotify client context
type key int

// CLIENT_CONTEXT is the key used for the spotify client context
var CLIENT_CONTEXT = key(1)

// spotifyClient interface of spotify client
type spotifyClient interface {
	CurrentUsersPlaylists() (*spotify.SimplePlaylistPage, error)
}

// playlistItem is a simplified structure of a playlist item
type playlistItem struct {
	Image     string      `json:"image"`
	Name      string      `json:"name"`
	OwnerName string      `json:"owner_name"`
	ID        spotify.ID  `json:"ID"`
	URI       spotify.URI `json:"uri"`
}

// reducePlaylist will reduce the spotify playlist to a simplified one
func reducePlaylist(playlistResp *spotify.SimplePlaylistPage) []playlistItem {
	playlist := []playlistItem{}

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

// playlistHandler is the handler to get the current user playlists
func playlistHandler(w http.ResponseWriter, r *http.Request) {
	client := r.Context().Value(CLIENT_CONTEXT).(spotifyClient)
	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		log.WithError(err).Error("playlistHandler: could not get user playlists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reducedPlaylist := reducePlaylist(playlists)
	json.NewEncoder(w).Encode(reducedPlaylist)
}

// tokenMiddleware will retrieve the token from the header and add the spotify client in the request context
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
	r.HandleFunc("/playlist", playlistHandler).Methods("GET")

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	contextedMux := tokenMiddleware(r)
	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(contextedMux)))
}
