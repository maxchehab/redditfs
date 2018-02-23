package main

import "github.com/maxchehab/geddit"

// Remove deletes a file from reddit
func Remove(file string, session *geddit.OAuthSession) {

}

// Create uploads a file from reddit
func Create(file string, session *geddit.OAuthSession) {
	session.Submit(geddit.NewTextSubmission("77346c3e708a", "title", "hello world", false, nil))

}

// Change edits a file from reddit
func Change(file string, session *geddit.OAuthSession) {
	session.EditUserText(geddit.NewEdit("this is an edit", "t3_7pni8t"))
}
