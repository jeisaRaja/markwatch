package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func watchFile(watcher *fsnotify.Watcher, path, tFname string) {
	defer watcher.Close()

	dir := filepath.Dir(path)
	log.Println("watchFile called", dir)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Println("modified file: ", event.Name)
					err := run(path, tFname, os.Stdout, false)
					if err != nil {
						log.Println("failed to re run")
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error: ", err)
			}
		}
	}()

	err := watcher.Add(dir)
	if err != nil {
		return
	}
	<-make(chan struct{})
}
