//go:build watch

package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	watchPaths := []string{"articles", "style.css", "main.go", "index.html", "article.html"}
	for _, path := range watchPaths {
		if err := watcher.Add(path); err != nil {
			log.Println("watch error:", err)
		}
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			fmt.Println("Changed:", event.Name)
			buildSite()
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("watch error:", err)
		}
	}
}
