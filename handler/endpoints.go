package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var root *string

func RegisterEndpoints(mux *http.ServeMux, version string, startTime time.Time) {
	mux.HandleFunc("/_/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"version": version,
			"uptime":  math.Round(time.Since(startTime).Seconds()*10) / 10,
			"root":    *root,
		})

	})
	mux.HandleFunc("/_/api/download/", handleDownload)
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
