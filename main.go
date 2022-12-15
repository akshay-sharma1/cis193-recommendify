package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"html/template"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserTopRead))
var state = "abc123"
var client *spotify.Client

var ctxt = context.Background()

func main() {
	// serve css and images
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/preferences", getPreferences)

	http.ListenAndServe(":8080", nil)
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", HomePage).Methods("GET")
	router.HandleFunc("/callback", completeAuth)
	return router
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("html/start.html")
	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	// redirect to authorization endpoint
	url := auth.AuthURL(state)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client = spotify.New(auth.Client(r.Context(), tok))
	http.Redirect(w, r, "/preferences", http.StatusSeeOther)
}

func getPreferences(w http.ResponseWriter, r *http.Request) {
	// get mood photos/images
	mood := GetMoodMetadata()
	top_tracks := getTopTrackMetadata(client, ctxt)
	preferences_data := PreferencesData{
		mood,
		top_tracks,
	}

	t, err := template.ParseFiles("html/preferences.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	if r.Method != http.MethodPost {
		t.Execute(w, preferences_data)
		return
	}
}
