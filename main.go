package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"html/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserTopRead, spotifyauth.ScopePlaylistModifyPublic))
var state = "abc123"
var client *spotify.Client

var (
	key   = []byte("recommendify-key")
	store = sessions.NewCookieStore(key)
)

var ctxt = context.Background()

func main() {
	// serve css and images
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/preferences", getPreferences)
	http.HandleFunc("/recommendations", getRecommendations)
	http.HandleFunc("/confirmation", getPlaylist)

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

	// add authenticated flag in session
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = true
	session.Save(r, w)

	// use the token to get an authenticated client
	client = spotify.New(auth.Client(r.Context(), tok))
	http.Redirect(w, r, "/preferences", http.StatusSeeOther)
}

func getPreferences(w http.ResponseWriter, r *http.Request) {
	// get mood photos/images
	mood := GetMoodMetadata()
	top_tracks := GetTopTrackMetadata(client, ctxt)
	genre_list := GenerateAutocomplete(client, ctxt)
	preferences_data := PreferencesData{
		mood,
		top_tracks,
		genre_list,
	}

	t, err := template.ParseFiles("html/preferences.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	if r.Method == http.MethodGet {
		t.Execute(w, preferences_data)
		return
	}

	if r.Method == http.MethodPost {
		session, _ := store.Get(r, "cookie-name")

		// get information from post request and redirect
		parse_err := r.ParseForm()
		if err != nil {
			fmt.Println(parse_err)
		}

		// parse form input
		if r.Form.Has("logout_input") {
			logout(w, r)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else if r.Form.Has("genreInput") {
			session.Values["genreInput"] = r.Form.Get("genreInput")
			session.Values["moodInput"] = nil
			session.Values["topTrackInput"] = nil
		} else if r.Form.Has("moodInput") {
			session.Values["moodInput"] = r.Form.Get("moodInput")
			session.Values["topTrackInput"] = nil
			session.Values["genreInput"] = nil
		} else {
			session.Values["topTrackInput"] = r.Form.Get("topTrackInput")
			session.Values["moodInput"] = nil
			session.Values["genreInput"] = nil
		}

		session.Save(r, w)
		http.Redirect(w, r, "/recommendations", http.StatusSeeOther)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Values["genreInput"] = nil
	session.Values["moodInput"] = nil
	session.Values["topTrackInput"] = nil
	session.Values["playlist_id"] = nil
	session.Values["track_ids"] = nil
	session.Save(r, w)
}

func getRecommendations(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, _ := store.Get(r, "cookie-name")

		t, err := template.ParseFiles("html/rec.html")
		if err != nil {
			fmt.Println(err)
			return
		}

		if session.Values["moodInput"] != nil {
			val, ok := session.Values["moodInput"].(string)
			if !ok {
				fmt.Println("moodInput isn't a string")
			}

			recommendations := RecommendMood(client, ctxt, strings.ToLower(val))
			session.Values["track_ids"] = recommendations.RecommendTrackID
			session.Save(r, w)
			t.Execute(w, recommendations)
		} else if session.Values["genreInput"] != nil {
			val, ok := session.Values["genreInput"].(string)
			if !ok {
				fmt.Println("genreInput isn't a string")
			}

			recommendations := RecommendFromGenre(client, ctxt, strings.ToLower(val))
			session.Values["track_ids"] = recommendations.RecommendTrackID
			session.Save(r, w)
			t.Execute(w, recommendations)
		} else if session.Values["topTrackInput"] != nil {
			val, ok := session.Values["topTrackInput"].(string)
			if !ok {
				fmt.Println("topTrackInput isn't a string")
			}

			recommendations := RecommendFromTrack(client, ctxt, val)
			session.Values["track_ids"] = recommendations.RecommendTrackID
			session.Save(r, w)
			t.Execute(w, recommendations)
		}
	}

	// if user selects logout or create playlist, handle that action
	if r.Method == http.MethodPost {
		session, _ := store.Get(r, "cookie-name")

		// get information from post request and redirect
		parse_err := r.ParseForm()
		if parse_err != nil {
			fmt.Println(parse_err)
		}

		if r.Form.Has("logout_input") {
			logout(w, r)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else if r.Form.Has("playlistInput") {
			val, ok := session.Values["track_ids"].([]string)
			if !ok {
				fmt.Println("error getting recommended track ids for playlist")
			}
			playlist_id := CreatePlaylist(client, ctxt, val)
			session.Values["playlist_id"] = playlist_id
		}
		session.Save(r, w)
		http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
	}
}

func getPlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		logout(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		session, _ := store.Get(r, "cookie-name")
		t, err := template.ParseFiles("html/confirmation.html")
		if err != nil {
			fmt.Println(err)
			return

		}
		if session.Values["playlist_id"] != nil {
			val, ok := session.Values["playlist_id"].(string)
			if !ok {
				fmt.Println("error parsing session playlist id")
			}
			playlist := GetPlaylist(client, ctxt, val)
			t.Execute(w, playlist)
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
