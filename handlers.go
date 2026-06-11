package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	metaLock    sync.RWMutex
	appMetadata = struct {
		Favorites map[string][]string `json:"favorites"`
		Recent    map[string][]string `json:"recent"`
	}{
		Favorites: make(map[string][]string),
		Recent:    make(map[string][]string),
	}
)

func saveMetadata() {
	metaLock.Lock()
	defer metaLock.Unlock()
	data, _ := json.Marshal(appMetadata)
	os.WriteFile("metadata.json", data, 0644)
}

func loadMetadata() {
	metaLock.Lock()
	defer metaLock.Unlock()
	data, err := os.ReadFile("metadata.json")
	if err == nil {
		json.Unmarshal(data, &appMetadata)
	}
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if strings.Contains(ip, ":") {
		host, _, err := net.SplitHostPort(ip)
		if err == nil {
			return host
		}
	}
	return ip
}

func addRecent(ip, path string) {
	if path == "" || strings.HasPrefix(path, "/.trash") {
		return
	}
	metaLock.Lock()
	defer metaLock.Unlock()
	list := appMetadata.Recent[ip]
	for i, p := range list {
		if p == path {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	list = append([]string{path}, list...)
	if len(list) > 20 {
		list = list[:20]
	}
	appMetadata.Recent[ip] = list
	go saveMetadata()
}

func isFavorite(ip, path string) bool {
	metaLock.RLock()
	defer metaLock.RUnlock()
	list := appMetadata.Favorites[ip]
	for _, p := range list {
		if p == path {
			return true
		}
	}
	return false
}

func storageHandler(w http.ResponseWriter, r *http.Request) {
	requestedPath := strings.TrimPrefix(r.URL.Path, "/storage")
	if requestedPath == "" {
		requestedPath = "/"
	}
	cleanPath := filepath.Clean(requestedPath)
	if strings.HasPrefix(cleanPath, "..") || cleanPath == ".." {
		http.Error(w, "Invalid path: Directory traversal attempt detected.", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}
	fullPath := filepath.Join(storageDir, cleanPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File or directory not found", http.StatusNotFound)
		} else {
			log.Printf("Error stating path %s: %v", fullPath, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	if info.IsDir() {
		serveDirectoryListing(w, r, fullPath, cleanPath)
	} else {
		addRecent(getIP(r), cleanPath)
		http.ServeFile(w, r, fullPath)
	}
}

func getDirSize(dirPath string) (int64, error) {
	var size int64
	err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func serveDirectoryListing(w http.ResponseWriter, r *http.Request, fullPath, currentPath string) {
	files, err := os.ReadDir(fullPath)
	if err != nil {
		log.Printf("Error reading directory %s: %v", fullPath, err)
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}
	var fileInfos []FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Printf("Error getting file info for %s: %v", filepath.Join(fullPath, file.Name()), err)
			continue
		}
		fPath := filepath.ToSlash(filepath.Join(currentPath, file.Name()))
		if file.Name() == ".trash" && currentPath == "/" {
			continue
		}
		fileInfos = append(fileInfos, FileInfo{
			Name:       file.Name(),
			Path:       fPath,
			IsDir:      file.IsDir(),
			Size:       info.Size(),
			ModTime:    info.ModTime().Format("2006-01-02 15:04:05"),
			IsFavorite: isFavorite(getIP(r), fPath),
		})
	}
	parentPath := ""
	if currentPath != "/" {
		parentPath = filepath.Dir(currentPath)
		if parentPath == "." {
			parentPath = "/"
		}
	}
	usedStorage, _ := getDirSize(storageDir)
	funcMap := template.FuncMap{
		"formatSize": func(bytes int64) string {
			const (
				kb = 1024
				mb = kb * 1024
				gb = mb * 1024
			)
			switch {
			case bytes < kb:
				return fmt.Sprintf("%d B", bytes)
			case bytes < mb:
				return fmt.Sprintf("%.2f KB", float64(bytes)/kb)
			case bytes < gb:
				return fmt.Sprintf("%.2f MB", float64(bytes)/mb)
			default:
				return fmt.Sprintf("%.2f GB", float64(bytes)/gb)
			}
		},
		"isImage": func(filename string) bool {
			ext := strings.ToLower(filepath.Ext(filename))
			return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp" || ext == ".webp"
		},
		"usagePercent": func(used, max int64) float64 {
			if max == 0 {
				return 0
			}
			p := (float64(used) / float64(max)) * 100
			if p > 100 {
				return 100
			}
			return p
		},
	}
	tmpl, err := template.New("fileList").Funcs(funcMap).Parse(fileListTemplate)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
	data := TemplateData{
		CurrentPath: currentPath,
		ParentPath:  parentPath,
		Files:       fileInfos,
		UsedStorage: usedStorage,
		MaxStorage:  maxStorageSize,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	targetDir := r.URL.Query().Get("dir")
	if targetDir == "" {
		targetDir = "/"
	}
	cleanTargetDir := filepath.Clean(targetDir)
	if strings.HasPrefix(cleanTargetDir, "..") || cleanTargetDir == ".." {
		http.Error(w, "Invalid target directory: Directory traversal attempt detected.", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(cleanTargetDir, "/") {
		cleanTargetDir = "/" + cleanTargetDir
	}
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}
	files := r.MultipartForm.File["uploadFile"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}
	currentUsed, _ := getDirSize(storageDir)
	for _, fileHeader := range files {
		if currentUsed+fileHeader.Size > maxStorageSize {
			log.Printf("Storage limit exceeded: cannot upload %s (%d bytes)", fileHeader.Filename, fileHeader.Size)
			continue
		}
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Error opening uploaded file %s: %v", fileHeader.Filename, err)
			continue
		}
		destPath := filepath.Join(storageDir, cleanTargetDir, fileHeader.Filename)
		dst, err := os.Create(destPath)
		if err != nil {
			log.Printf("Error creating destination file %s: %v", destPath, err)
			file.Close()
			continue
		}
		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("Error copying file content to %s: %v", destPath, err)
			dst.Close()
			file.Close()
			continue
		}
		addRecent(getIP(r), filepath.ToSlash(filepath.Join(cleanTargetDir, fileHeader.Filename)))
		currentUsed += fileHeader.Size
		dst.Close()
		file.Close()
	}
	http.Redirect(w, r, "/storage"+cleanTargetDir, http.StatusSeeOther)
}

func mkdirHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	targetDir := r.FormValue("dir")
	folderName := r.FormValue("name")
	if folderName == "" {
		http.Error(w, "Folder name is required", http.StatusBadRequest)
		return
	}
	cleanTargetDir := filepath.Clean(targetDir)
	fullPath := filepath.Join(storageDir, cleanTargetDir, folderName)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		log.Printf("Error creating directory %s: %v", fullPath, err)
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/storage"+cleanTargetDir, http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	itemPath := r.FormValue("path")
	cleanPath := filepath.Clean(itemPath)
	if strings.HasPrefix(cleanPath, "..") || cleanPath == ".." {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	fullPath := filepath.Join(storageDir, cleanPath)
	if strings.HasPrefix(cleanPath, "/.trash") {
		if err := os.RemoveAll(fullPath); err != nil {
			log.Printf("Error deleting %s: %v", fullPath, err)
			http.Error(w, "Error deleting item", http.StatusInternalServerError)
			return
		}
	} else {
		trashPath := filepath.Join(storageDir, ".trash", filepath.Base(cleanPath))
		if _, err := os.Stat(trashPath); err == nil {
			trashPath = filepath.Join(storageDir, ".trash", fmt.Sprintf("%d_%s", os.Getpid(), filepath.Base(cleanPath)))
		}
		if err := os.Rename(fullPath, trashPath); err != nil {
			log.Printf("Error trashing %s: %v", fullPath, err)
			http.Error(w, "Error moving to trash", http.StatusInternalServerError)
			return
		}
	}

	parentDir := filepath.Dir(cleanPath)
	if !strings.HasPrefix(parentDir, "/") {
		parentDir = "/" + parentDir
	}
	http.Redirect(w, r, "/storage"+filepath.ToSlash(parentDir), http.StatusSeeOther)
}

func favoriteToggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := r.FormValue("path")
	ip := getIP(r)
	metaLock.Lock()
	list := appMetadata.Favorites[ip]
	found := -1
	for i, p := range list {
		if p == path {
			found = i
			break
		}
	}
	if found != -1 {
		appMetadata.Favorites[ip] = append(list[:found], list[found+1:]...)
	} else {
		appMetadata.Favorites[ip] = append(list, path)
	}
	metaLock.Unlock()
	saveMetadata()
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func favoritesHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	metaLock.RLock()
	paths := appMetadata.Favorites[ip]
	metaLock.RUnlock()
	serveVirtualListing(w, r, "Favorites", paths)
}

func recentHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	metaLock.RLock()
	paths := appMetadata.Recent[ip]
	metaLock.RUnlock()
	serveVirtualListing(w, r, "Recent", paths)
}

func trashHandler(w http.ResponseWriter, r *http.Request) {
	serveDirectoryListing(w, r, filepath.Join(storageDir, ".trash"), ".trash")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/storage/", http.StatusSeeOther)
		return
	}

	var paths []string
	filepath.Walk(storageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(storageDir, path)
		fPath := "/" + filepath.ToSlash(rel)
		if fPath == "/" || strings.HasPrefix(fPath, "/.trash") || info.Name() == "metadata.json" {
			return nil
		}
		if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(query)) {
			paths = append(paths, fPath)
		}
		return nil
	})
	serveVirtualListing(w, r, "Search: "+query, paths)
}

func serveVirtualListing(w http.ResponseWriter, r *http.Request, title string, paths []string) {
	var fileInfos []FileInfo
	for _, p := range paths {
		fullPath := filepath.Join(storageDir, p)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, FileInfo{
			Name:       filepath.Base(p),
			Path:       p,
			IsDir:      info.IsDir(),
			Size:       info.Size(),
			ModTime:    info.ModTime().Format("2006-01-02 15:04:05"),
			IsFavorite: isFavorite(getIP(r), p),
		})
	}
	usedStorage, _ := getDirSize(storageDir)
	data := TemplateData{
		CurrentPath: title,
		ParentPath:  "",
		Files:       fileInfos,
		UsedStorage: usedStorage,
		MaxStorage:  maxStorageSize,
	}
	tmpl, _ := template.New("fileList").Funcs(template.FuncMap{
		"formatSize":   func(b int64) string { return fmt.Sprintf("%.2f MB", float64(b)/(1024*1024)) },
		"isImage":      func(n string) bool { ext := strings.ToLower(filepath.Ext(n)); return ext == ".jpg" || ext == ".png" },
		"usagePercent": func(u, m int64) float64 { return (float64(u) / float64(m)) * 100 },
	}).Parse(fileListTemplate)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}
