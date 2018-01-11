package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
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
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("create file:", event.Name)
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("changed file: ", event.Name)
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./test")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
