
# Recommendify

[![Go Report Card](https://goreportcard.com/badge/github.com/akshay-sharma1/cis193-recommendify)](https://goreportcard.com/report/github.com/akshay-sharma1/cis193-recommendify)

> Discover new music!

![](landingPage.PNG)

Recommendify is a web application that allows Spotify users to get personalized music recommendations based on a particular genre, mood, or favorite track. 

The algorithm works by matching a user's particular input with a "seed" catalog of their top artists and tracks, in order to 
curate only the most relevant tracks. You can then preview individual songs, add the playlist to your Spotify library, and start listening to it right away!

## Usage Instructions
To use this app, first set the following environment variables:

``export SPOTIFY_SECRET=2bbb563a76344f5093984312db7c7a1f``

``export SPOTIFY_ID=ed2847982e1741148c1e48512bc5fa55``

Then, run the following command: 

``go run main.go api.go`` 

and navigate to ``localhost:8080`` in your browser to test the web application out.


## Spotify API
The application consumes the following endpoints of the Spotify API:
 * [User Authorization](https://developer.spotify.com/documentation/general/guides/authorization-guide/)
 * [Search for an Item](https://developer.spotify.com/documentation/web-api/reference/search/search/)
 * [Get User Top Tracks/Artists](https://developer.spotify.com/documentation/web-api/reference/personalization/get-users-top-artists-and-tracks/)
 * [Seed Recommendations](https://developer.spotify.com/web-api/get-recommendations/)
 * [Create Playlist](https://developer.spotify.com/documentation/web-api/reference/playlists/create-playlist/)
 * [Add Tracks to a Playlist](https://developer.spotify.com/documentation/web-api/reference/playlists/add-tracks-to-playlist/)
 

## Stack
* Golang
* HTML/CSS/JS


##  License
MIT

