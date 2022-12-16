package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/zmb3/spotify/v2"
)

type PreferencesData struct {
	Mood
	TopTracks
	Genre
}

type Mood struct {
	Identifier []string
	Image      []string
}

type TopTracks struct {
	Name       []string
	SongImage  []string
	ArtistName []string
}

type Genre struct {
	Genres []string
}

func GetMoodMetadata() Mood {
	BASE_URL := "https://source.unsplash.com/"

	image_ids := []string{"dWIVg59BVXY", "fnztlIb52gU", "s9CC2SKySJM", "zfPOelmDc-M"}
	image_urls := make([]string, 4)

	for i := 0; i < 4; i++ {
		image_urls[i] = BASE_URL + image_ids[i]
	}

	identifier_list := []string{"Chill", "Mood Booster", "Deep Focus", "Workout"}

	new_mood := Mood{
		Identifier: identifier_list,
		Image:      image_urls,
	}

	return new_mood
}

func GetTopTrackMetadata(client *spotify.Client, ctxt context.Context) TopTracks {
	tracks, err := client.CurrentUsersTopTracks(ctxt, spotify.Limit(10))
	if err != nil {
		fmt.Println(err)
		return TopTracks{}
	}

	song_names := make([]string, 10)
	song_images := make([]string, 10)
	artist_names := make([]string, 10)

	for i, elem := range tracks.Tracks {
		song_names[i] = elem.Name

		for j, image := range elem.Album.Images {
			if j == 1 {
				song_images[i] = image.URL
			}
		}

		for k, artist := range elem.Artists {
			if k == 0 {
				artist_names[i] = artist.Name
			}
		}
	}

	return TopTracks{
		Name:       song_names,
		SongImage:  song_images,
		ArtistName: artist_names,
	}
}

func GenerateAutocomplete(client *spotify.Client, ctxt context.Context) Genre {
	genre_list, err := client.GetAvailableGenreSeeds(ctxt)
	if err != nil {
		fmt.Println(err)
	}

	sort.Strings(genre_list)
	return Genre{
		genre_list,
	}
}
