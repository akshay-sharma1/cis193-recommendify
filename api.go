package main

import "github.com/zmb3/spotify/v2"

type Mood struct {
	Identifier []string
	Image      []string
}

type TopTracks struct {
}

type TopArtists struct {
}

func GetMoodMetadata() Mood {
	// TODO: change image links to proper source
	BASE_URL := "https://source.unsplash.com"

	image_ids := []string{"dWIVg59BVXY", "fnztlIb52gU", "s9CC2SKySJM", "zfPOelmDc-M"}
	image_urls := make([]string, 4)

	for i := 0; i < 4; i++ {
		image_urls[i] = BASE_URL + image_ids[i]
	}

	identifier_list := []string{"Chill", "Mood Booster", "Deep Focus", "Workout"}
	identifier_map := make(map[string]string)

	for i := 0; i < 4; i++ {
		identifier_map[identifier_list[i]] = image_urls[i]
	}

	new_mood := Mood{
		Identifier: identifier_list,
		Image:      image_urls,
	}

	return new_mood
}

func getTopTrackMetadata(client *spotify.Client) {

}
