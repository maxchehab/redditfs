package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Manifest data structure describing the location of a file
// {
// 	"repositories":[
//    {
// 	  "name":"My dankeroni repository",
// 	  "files":[
// 		 {
// 			"name":"test.txt",
// 			"path":"/",
// 			"location":[
// 			   "7zp1o6"
// 			]
// 		 }
// 	  ]
//    }
// 	]
// }
type Manifest struct {
	Repositories []Repository `json:"repositories"`
}

// Repository structure
type Repository struct {
	Name  string `json:"name"`
	Files []File `json:"files"`
}

// File structure
type File struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Location []string `json:"location"`
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

// FilterReposByNames filters repositories with a collection of names.
func (m Manifest) FilterReposByNames(responses []string) (repositories []Repository) {
	for _, repo := range m.Repositories {
		for _, response := range responses {
			if repo.Name == response {
				repositories = append(repositories, repo)
			}
		}
	}
	return
}

// Download a file
func (f File) Download(path string) (err error) {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}
	for _, location := range f.Location {
		url := fmt.Sprintf(`https://www.reddit.com/r/77346c3e708a/comments/%v.json`, location)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		request.Header.Set("User-Agent", userAgent)
		client := &http.Client{}

		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf(`bad status code, [%v]`, response.StatusCode)
		}

		JSON, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		listing, err := CreateListingFromByteArray(JSON)
		err = WriteByteStringToFile(listing.Text, path)
		if err != nil {
			return err
		}
	}
	return
}

// Download a repository
func (r Repository) Download(path string) (err error) {
	for _, file := range r.Files {
		fmt.Printf("Downloading /%v%v%v\n", r.Name, file.Path, file.Name)
		file.Download(path + "/" + r.Name + file.Path + file.Name)
	}
	return
}

// RetrieveManifestFromReddit will download a manifest from a specified subreddit
func RetrieveManifestFromReddit(subreddit string) (Manifest, error) {
	// https://www.reddit.com/r/[repo]/search.json?q=manifest.json&restrict_sr=on&sort=relevance&t=all
	var m Manifest
	url := fmt.Sprintf(`https://www.reddit.com/r/%v/search.json?q=manifest&restrict_sr=on&sort=relevance&t=all`, subreddit)

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
