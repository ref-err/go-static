package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var root *string

func RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/_/status", handleStatus)
	mux.HandleFunc("/_/api/download/", handleDownload)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	requestedPath := strings.TrimPrefix(r.URL.Path, "/_/api/download/")

	rootDir, err := filepath.Abs(*root)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	safePath := filepath.Join(rootDir, filepath.Clean(requestedPath))

	if !strings.HasPrefix(safePath, rootDir) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	if _, err := os.Stat(safePath); os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(safePath)))
	http.ServeFile(w, r, safePath)
}

func SetRoot(value *string) {
	root = value
}
