package main

func getMoodMetadata() map[string]string {
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

	return identifier_map
}
