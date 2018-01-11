package main

import (
	"log"

	"github.com/maxchehab/geddit"
)

// GetSession creates geddit session using OAuth
func GetSession() *geddit.OAuthSession {
	session, err := geddit.NewOAuthSession(
		"gNl1rziyJUjwNQ",
		"TdeaiSX6FaxwBsfpi9L6FxFu288",
		"redditfs",
		"http://maxchehab.com",
	)

	if err != nil {
		log.Fatal(err)
	}

	username, password := Credentials()

	// Create new auth token for confidential clients (personal scripts/apps).
	err = session.LoginAuth(username, password)
	if err != nil {
		log.Fatal(err)
	}

	return session
}
