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

// User is a simplified structure of a user
type User struct {
	Name  string `json:'name'`
	ID    string `json:'id'`
	Image string `json:'image'`
}

// userHandler is the handler to get the current user info
func userHandler(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := new(oauth2.Token)
	token.AccessToken = bearer
	client := spotify.Authenticator{}.NewClient(token)
	user, err := client.CurrentUser()
	if err != nil {
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

// userFromHandler is the handler to get user info from his spotify ID
func userFromHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userID"]

	bearer := r.Header.Get("Authorization")
	token := new(oauth2.Token)
	token.AccessToken = bearer
	client := spotify.Authenticator{}.NewClient(token)
	user, err := client.GetUsersPublicProfile(spotify.ID(userId))
	if err != nil {
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/user", userHandler).Methods("GET")
	r.HandleFunc("/user/{userID}", userFromHandler).Methods("GET")

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(r)))
}
