package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	serverPort     = "50"
	wwwDir         = "./www"
	storageDir     = "./storage"
	maxStorageSize = 1000 * 1024 * 1024 * 1024 // Configurable limit (e.g., 1 GB)
)

func main() {
	if err := os.MkdirAll(wwwDir, 0755); err != nil {
		log.Fatalf("Failed to create www directory: %v", err)
	}
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(storageDir, ".trash"), 0755); err != nil {
		log.Fatalf("Failed to create trash directory: %v", err)
	}
	loadMetadata()
	http.Handle("/", http.FileServer(http.Dir(wwwDir)))
	http.HandleFunc("/storage/", storageHandler)
	http.HandleFunc("/storage/upload", uploadHandler)
	http.HandleFunc("/storage/mkdir", mkdirHandler)
	http.HandleFunc("/storage/delete", deleteHandler)
	http.HandleFunc("/storage/favorite/toggle", favoriteToggleHandler)
	http.HandleFunc("/storage/favorites", favoritesHandler)
	http.HandleFunc("/storage/recent", recentHandler)
	http.HandleFunc("/storage/trash/", trashHandler)
	http.HandleFunc("/storage/search", searchHandler)
	addr := fmt.Sprintf(":%s", serverPort)
	log.Printf("Go NAS server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
