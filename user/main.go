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
	GetUsersPublicProfile(userID spotify.ID) (*spotify.User, error)
	CurrentUser() (*spotify.PrivateUser, error)
}

// User is a simplified structure of a user
type User struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Image string `json:"image"`
}

// userHandler is the handler to get the current user info
func userHandler(w http.ResponseWriter, r *http.Request) {
	client := r.Context().Value(CLIENT_CONTEXT).(spotifyClient)
	user, err := client.CurrentUser()
	if err != nil {
		log.WithError(err).Error("userHandler: could not get current user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Send critical info to other microservice to keep track of country / product / date of birth

	var image string
	if len(user.Images) > 0 {
		image = user.Images[0].URL
	}

	simplifiedUser := User{
		Name:  user.DisplayName,
		ID:    user.ID,
		Image: image,
	}
	json.NewEncoder(w).Encode(simplifiedUser)
}

// userFromHandler is the handler to get user info from his spotify ID
func userFromHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userID"]
	client := r.Context().Value(CLIENT_CONTEXT).(spotifyClient)
	user, err := client.GetUsersPublicProfile(spotify.ID(userId))
	if err != nil {
		log.WithField("userID", userId).WithError(err).Error("userFromHandler: could not get user public profile")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var image string
	if len(user.Images) > 0 {
		image = user.Images[0].URL
	}

	simplifiedUser := User{
		Name:  user.DisplayName,
		ID:    user.ID,
		Image: image,
	}
	json.NewEncoder(w).Encode(simplifiedUser)
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
	r.HandleFunc("/user", userHandler).Methods("GET")
	r.HandleFunc("/user/{userID}", userFromHandler).Methods("GET")

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Content-Type", "Accept", "*"},
	})

	contextedMux := tokenMiddleware(r)
	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(contextedMux)))
}
