package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/maxchehab/geddit"
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

// Contains checks if a manifest contians selected repository
func (m *Manifest) Contains(repository Repository) bool {
	for _, repo := range m.Repositories {
		if repo.Name == repository.Name {
			return true
		}
	}
	return false
}

// UploadBuffer uploads a buffer of data and modifies the file object
func (f *File) UploadBuffer(buffer []byte, session *geddit.OAuthSession) (err error) {
	text := ""
	for _, b := range buffer {
		o := strconv.Itoa(int(b))
		text += o + " "
	}
	submission, err := session.Submit(geddit.NewTextSubmission(testSubreddit, testSubreddit, text, false, nil))
	if err != nil {
		return
	}
	f.Location = append(f.Location, submission.ID)
	return
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
func (f File) Download(path string, session *geddit.OAuthSession) (err error) {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}
	for _, location := range f.Location {
		url := fmt.Sprintf(`https://oauth.reddit.com/r/77346c3e708a/comments/%v.json`, location)

		JSON, err := session.GetRawRequest(url)
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
func (r Repository) Download(path string, session *geddit.OAuthSession) (err error) {
	for _, file := range r.Files {
		fmt.Printf("Downloading /%v%v%v\n", r.Name, file.Path, file.Name)
		file.Download(path+"/"+r.Name+file.Path+file.Name, session)
	}
	return
}

// RetrieveManifestFromReddit will download a manifest from a specified subreddit
func RetrieveManifestFromReddit(subreddit string, session *geddit.OAuthSession) (m Manifest, err error) {
	// https://oauth.reddit.com/r/[repo]/search.json?q=manifest.json&restrict_sr=on&sort=relevance&t=all
	url := fmt.Sprintf(`https://oauth.reddit.com/r/%v/search.json?q=manifest&restrict_sr=on&sort=relevance&t=all`, subreddit)

	JSON, err := session.GetRawRequest(url)
	if err != nil {
		return
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

// package main

// import "github.com/maxchehab/geddit"

// // Remove deletes a file from reddit
// func Remove(file string, session *geddit.OAuthSession) {

// }

// // Upload uploads a file from reddit
// func Upload(file string, session *geddit.OAuthSession) {
// 	session.Submit(geddit.NewTextSubmission("77346c3e708a", "title", "hello world", false, nil))

// }

// // Change edits a file from reddit
// func Change(file string, session *geddit.OAuthSession) {
// 	session.EditUserText(geddit.NewEdit("this is an edit", "t3_7pni8t"))
// }
