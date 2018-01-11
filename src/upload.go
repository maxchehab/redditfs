package main

import "github.com/maxchehab/geddit"

func remove(file string, session *geddit.OAuthSession) {

}

func create(file string, session *geddit.OAuthSession) {
	session.Submit(geddit.NewTextSubmission("77346c3e708a", "title", "hello world", false, nil))

}

func change(file string, session *geddit.OAuthSession) {
	session.EditUserText(geddit.NewEdit("this is an edit", "t3_7pni8t"))
}
