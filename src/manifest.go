package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Manifest data structure describing the location of a file
// {
// 	"files": [
// 	  {
// 		"name": "test.txt",
// 		"path": [
// 		  "7zp1o6"
// 		]
// 	  }
// 	]
// }
type Manifest struct {
	Files []struct {
		Name string   `json:"name"`
		Path []string `json:"path"`
	} `json:"files"`
}

// CreateManifestFromByteArray creates a Manifest object from a byte array
func CreateManifestFromByteArray(JSON []byte) (Manifest, error) {
	var m Manifest
	err := json.Unmarshal(JSON, &m)

	return m, err
}

// CreateManifestFromString creates a manifest object from a JSON string
func CreateManifestFromString(JSON string) (Manifest, error) {
	var m Manifest
	err := json.Unmarshal([]byte(JSON), &m)

	return m, err
}

// RetrieveManifestFromReddit will download a manifest from a specified subreddit
func RetrieveManifestFromReddit(subreddit string) (Manifest, error) {
	// https://www.reddit.com/r/[repo]/search.json?q=manifest.json&restrict_sr=on&sort=relevance&t=all
	var m Manifest
	url := fmt.Sprintf(`https://www.reddit.com/r/%v/search.json?q=manifest&restrict_sr=on&sort=relevance&t=all`, subreddit)

	// response, err := client.Get(request)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return m, err
	}
	request.Header.Set("User-Agent", userAgent)
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return m, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println(request)
		return m, fmt.Errorf(`bad status code, [%v]`, response.StatusCode)
	}

	JSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return m, err
	}

	search, err := CreateSearchFromByteArray(JSON)
	if len(search.Data.Children) == 0 {
		return m, errors.New("could not locate manifest.json")
	}

	for _, listing := range search.Data.Children {
		if listing.Data.Subreddit == subreddit && listing.Data.Title == "manifest.json" {
			return CreateManifestFromString(listing.Data.Text)
		}
	}

	return m, errors.New("could not locate manifest.json")
}
