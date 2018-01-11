package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/maxchehab/geddit"
)

func watch(path string, session *geddit.OAuthSession) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("remove file:", event.Name)
					watcher.Remove(event.Name)
					remove(event.Name, session)
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("create file:", event.Name)
					watcher.Add(event.Name)
					create(event.Name, session)
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("changed file: ", event.Name)
					change(event.Name, session)
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	folders := Subfolders(path)
	for _, folder := range folders {
		err = watcher.Add(folder)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
