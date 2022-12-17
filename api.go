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
	SongIds    []string
}

type Genre struct {
	Genres []string
}

type Recommendations struct {
	RecommendName    []string
	RecommendImage   []string
	RecommendArtist  []string
	RecommendSpotify []string
	RecommendTrackID []string
}

type CreatedPlaylist struct {
	PlaylistImg     string
	PlaylistName    string
	PlaylistPreview string
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

func getAlbumImageURI(client *spotify.Client, ctxt context.Context, id spotify.ID) spotify.Image {
	tracks, err := client.GetTrack(ctxt, id)
	if err != nil {
		fmt.Println("error getting track: ", err)
	}

	for i, elem := range tracks.Album.Images {
		if i == 0 {
			return elem
		}
	}

	return spotify.Image{}
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
	song_ids := make([]string, 10)

	for i, elem := range tracks.Tracks {
		song_names[i] = elem.Name
		song_ids[i] = string(elem.ID)

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
		SongIds:    song_ids,
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

func RecommendMood(client *spotify.Client, ctxt context.Context, mood string) Recommendations {
	var tgtValence, tgtEnergy float64
	targetGenres := make([]string, 2)
	if mood == "chill" {
		tgtValence = 0.6
		tgtEnergy = 0.2
		targetGenres = []string{"chill", "pop"}
	} else if mood == "mood booster" {
		tgtValence = 0.8
		tgtEnergy = 0.5
		targetGenres = []string{"pop", "happy"}
	} else if mood == "deep focus" {
		tgtValence = 0.6
		tgtEnergy = 0.2
		targetGenres = []string{"study", "classical"}
	} else {
		tgtValence = 0.7
		tgtEnergy = 0.8
		targetGenres = []string{"work-out", "hip-hop"}
	}

	seeds := spotify.Seeds{
		Genres: targetGenres,
	}

	track_attributes := spotify.NewTrackAttributes().TargetValence(tgtValence).TargetEnergy(tgtEnergy)

	recommendations, err := client.GetRecommendations(ctxt, seeds, track_attributes, spotify.Limit(30))
	if err != nil {
		fmt.Println("error getting recommendations: ", err)
		return Recommendations{}
	}

	song_names := make([]string, 30)
	artist_names := make([]string, 30)
	song_imgs := make([]string, 30)
	preview_urls := make([]string, 30)
	track_ids := make([]string, 30)

	for i, elem := range recommendations.Tracks {
		song_names[i] = elem.Name
		preview_urls[i] = elem.PreviewURL
		track_ids[i] = string(elem.ID)

		song_imgs[i] = getAlbumImageURI(client, ctxt, elem.ID).URL

		for k, artist := range elem.Artists {
			if k == 0 {
				artist_names[i] = artist.Name
			}
		}
	}

	return Recommendations{
		RecommendName:    song_names,
		RecommendImage:   song_imgs,
		RecommendArtist:  artist_names,
		RecommendSpotify: preview_urls,
		RecommendTrackID: track_ids,
	}
}

func RecommendFromGenre(client *spotify.Client, ctxt context.Context, genre string) Recommendations {
	genres := make([]string, 1)
	genres[0] = genre
	seeds := spotify.Seeds{
		Genres: genres,
	}

	recommendations, err := client.GetRecommendations(ctxt, seeds, nil, spotify.Limit(30))
	if err != nil {
		fmt.Println("error getting recommendations: ", err)
		return Recommendations{}
	}

	song_names := make([]string, 30)
	artist_names := make([]string, 30)
	song_imgs := make([]string, 30)
	preview_urls := make([]string, 30)
	track_ids := make([]string, 30)

	for i, elem := range recommendations.Tracks {
		song_names[i] = elem.Name
		preview_urls[i] = elem.PreviewURL
		track_ids[i] = string(elem.ID)

		song_imgs[i] = getAlbumImageURI(client, ctxt, elem.ID).URL

		for k, artist := range elem.Artists {
			if k == 0 {
				artist_names[i] = artist.Name
			}
		}
	}

	return Recommendations{
		RecommendName:    song_names,
		RecommendImage:   song_imgs,
		RecommendArtist:  artist_names,
		RecommendSpotify: preview_urls,
		RecommendTrackID: track_ids,
	}
}

func RecommendFromTrack(client *spotify.Client, ctxt context.Context, song_id string) Recommendations {
	song_ids := []spotify.ID{spotify.ID(song_id)}
	seeds := spotify.Seeds{
		Tracks: song_ids,
	}

	recommendations, err := client.GetRecommendations(ctxt, seeds, nil, spotify.Limit(30))
	if err != nil {
		fmt.Println("error getting recommendations: ", err)
		return Recommendations{}
	}

	song_names := make([]string, 30)
	artist_names := make([]string, 30)
	song_imgs := make([]string, 30)
	preview_urls := make([]string, 30)
	track_ids := make([]string, 30)

	for i, elem := range recommendations.Tracks {
		song_names[i] = elem.Name
		preview_urls[i] = elem.PreviewURL
		track_ids[i] = string(elem.ID)

		song_imgs[i] = getAlbumImageURI(client, ctxt, elem.ID).URL

		for k, artist := range elem.Artists {
			if k == 0 {
				artist_names[i] = artist.Name
			}
		}
	}

	return Recommendations{
		RecommendName:    song_names,
		RecommendImage:   song_imgs,
		RecommendArtist:  artist_names,
		RecommendSpotify: preview_urls,
		RecommendTrackID: track_ids,
	}
}

func CreatePlaylist(client *spotify.Client, ctxt context.Context, track_ids []string) string {
	curr_user, err := client.CurrentUser(ctxt)
	if err != nil {
		fmt.Println("error getting current user for creating playlist")
	}

	playlist, err := client.CreatePlaylistForUser(ctxt, curr_user.ID, "Your Recommended Playlist", "recommendify-generated playlist", true, false)
	if err != nil {
		fmt.Println("error creating playlist")
	}

	song_ids := make([]spotify.ID, len(track_ids))
	for i, ele := range track_ids {
		song_ids[i] = spotify.ID(ele)
	}
	client.AddTracksToPlaylist(ctxt, playlist.ID, song_ids...)

	return string(playlist.ID)
}

func GetPlaylist(client *spotify.Client, ctxt context.Context, playlist_id string) CreatedPlaylist {
	playlist, err := client.GetPlaylist(ctxt, spotify.ID(playlist_id))
	if err != nil {
		fmt.Println("error getting playlist for display")
		return CreatedPlaylist{}
	}

	var playlist_img, playlist_url, playlist_name string
	playlist_name = playlist.Name
	playlist_url = playlist.ExternalURLs["spotify"]

	for i, elem := range playlist.Images {
		if i == 0 {
			playlist_img = elem.URL
			break
		}
	}

	return CreatedPlaylist{
		PlaylistImg:     playlist_img,
		PlaylistName:    playlist_name,
		PlaylistPreview: playlist_url,
	}
}
