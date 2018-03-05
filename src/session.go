package main

import (
	"log"
	"os"
	"os/user"

	"github.com/maxchehab/geddit"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// GetSession creates geddit session using OAuth
func GetSession(username string, password string) (session *geddit.OAuthSession, err error) {
	// session, err := geddit.NewOAuthSession(
	// 	"EV4dAIStivu3XA",
	// 	"KnLTb7s6CrP9KlX_vhLwQ27VMW4",
	// 	"redditfs",
	// 	"http://maxchehab.com",
	//	77346c3e708a
	// )

	name, id, secret, read, err := GetOauthCredentials()
	if err != nil {
		return
	}

	session, err = geddit.NewOAuthSession(
		id,
		secret,
		name,
		"",
	)

	if err != nil {
		log.Fatal(err)
	}

	// Create new auth token for confidential clients (personal scripts/apps).
	err = session.LoginAuth(username, password)
	if err != nil {
		if err.Error() == "oauth2: cannot fetch token: 401 Unauthorized\nResponse: {\"message\": \"Unauthorized\", \"error\": 401}" && read {
			removeFile := false
			removeFilePrompt := &survey.Confirm{
				Message: "The OAUTH credentials provided for your authorized reddit application were incorrect. Would you like to change your credentials?",
			}
			survey.AskOne(removeFilePrompt, &removeFile, nil)

			if removeFile {

				usr, err := user.Current()
				if err != nil {
					return nil, err
				}

				credentialPath := usr.HomeDir + "/.redditfscredentials"
				os.Remove(credentialPath)
				return GetSession(username, password)
			}
		} else {
			return nil, err
		}
	}

	return
}
