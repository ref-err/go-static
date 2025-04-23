package handler

import (
	_ "embed"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"
)

// HTML template (renders on file server)

//go:embed templates/index.html
var indexTemplate string

// FileInfo contains basic info about a file.
type FileInfo struct {
	Name  string  // file name
	Path  string  // file path
	Size  float64 // file size
	IsDir bool    // is file a directory
}

// This is a re-implementation of default http.FileServer handler, written to render custom HTML/CSS.
func FileServer(root http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := filepath.ToSlash(r.URL.Path)

		f, err := root.Open(cleanPath)
		if err != nil { // error checking
			if os.IsPermission(err) {
				http.Error(w, "Access Denied", http.StatusForbidden)
				return
			}
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		defer f.Close()

		stat, _ := f.Stat() // gettings file metadata
		if !stat.IsDir() {
			http.FileServer(root).ServeHTTP(w, r) // serve file if its not a directory
			return
		}

		files, err := f.Readdir(-1) // read directory and get files
		if err != nil {
			http.Error(w, "Error reading dir:", http.StatusInternalServerError)
			return
		}

		fileList := make([]FileInfo, 0)
		for _, file := range files { // going thru every file
			fileList = append(fileList, FileInfo{ // adding file to fileList
				Name:  file.Name(),
				Path:  filepath.Join(cleanPath, file.Name()),
				Size:  math.Round((float64(file.Size())/1024.0)*10) / 10,
				IsDir: file.IsDir(),
			})
		}

		// creating a template out of its content
		tmpl, err := template.New("dir").Parse(indexTemplate)
		if err != nil {
			http.Error(w, "Template Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// rendering template
		tmpl.Execute(w, struct {
			Dir   string
			Files []FileInfo
		}{
			Dir:   cleanPath,
			Files: fileList,
		})
	})
}
