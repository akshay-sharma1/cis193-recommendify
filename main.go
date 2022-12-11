package main

import (
	"fmt"
	"log"
	"net/http"

	"html/template"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
var ch = make(chan *spotify.Client)
var state = "abc123"

func main() {
	// init router

	// serve css and images
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/callback", completeAuth)

	http.ListenAndServe(":8080", nil)
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", HomePage).Methods("GET")
	router.HandleFunc("/callback", completeAuth)
	return router
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("html/start.html")
	if err != nil {
		fmt.Println("template parsing error: ", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		fmt.Println("template executing error:", err)
	}
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
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
