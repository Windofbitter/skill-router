package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/dist/*
var webFS embed.FS

func getFileSystem() http.FileSystem {
	subFS, _ := fs.Sub(webFS, "web/dist")
	return http.FS(subFS)
}
