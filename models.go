package main

type FileInfo struct {
	Name       string
	Path       string
	IsDir      bool
	Size       int64
	ModTime    string
	IsFavorite bool
}

type TemplateData struct {
	CurrentPath string
	ParentPath  string
	Files       []FileInfo
	UsedStorage int64
	MaxStorage  int64
}
