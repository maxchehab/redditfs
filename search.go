package main

import (
	"encoding/json"
	"errors"
)

// Search data structure describing the results of a subreddit search
// {
// 	"data": {
// 	  "children": [
// 		{
// 		  "kind": "t3",
// 		  "data": {
// 			"subreddit": "77346c3e708a",
// 			"selftext": "test",
// 			"id": "7zpvqn",
// 			"author": "Senior-Jesticle",
// 			"score": 1,
// 			"name": "t3_7zpvqn",
// 			"url": "https://www.reddit.com/r/77346c3e708a/comments/7zpvqn/manifestjson/",
// 			"title": "manifest.json"
// 		  }
// 		}
// 	  ]
// 	}
// }
type Search struct {
	Data struct {
		Children []struct {
			Kind string  `json:"kind"`
			Data Listing `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// Listing structure
type Listing struct {
	Subreddit string `json:"subreddit"`
	Text      string `json:"selftext"`
	ID        string `json:"id"`
	Author    string `json:"author"`
	Score     int    `json:"score"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Title     string `json:"title"`
}

// CreateSearchFromByteArray creates a Search object from a byte array
func CreateSearchFromByteArray(JSON []byte) (Search, error) {
	var s Search
	err := json.Unmarshal(JSON, &s)

	return s, err
}

// CreateSearchFromJSON creates a Search object from a JSON string
func CreateSearchFromJSON(JSON string) (Search, error) {
	var s Search
	err := json.Unmarshal([]byte(JSON), &s)

	return s, err
}

// CreateListingFromByteArray creates a Listing object from a byte array
func CreateListingFromByteArray(JSON []byte) (Listing, error) {
	var s []Search
	err := json.Unmarshal(JSON, &s)
	if len(s) == 0 {
		return Listing{}, errors.New("could not find listing")
	}
	if len(s[0].Data.Children) == 0 {
		return Listing{}, errors.New("could not find listing")
	}
	return s[0].Data.Children[0].Data, err
}

// CreateListingFromJSON creates a Listing object from a JSON string
func CreateListingFromJSON(JSON string) (Listing, error) {
	var s Search
	err := json.Unmarshal([]byte(JSON), &s)
	if len(s.Data.Children) == 0 {
		return Listing{}, errors.New("could not find listing")
	}
	return s.Data.Children[0].Data, err
}
