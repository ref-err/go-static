package handler

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// HTML template (renders on file server)
const dirTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Index of {{.Dir}}</title>
	<style>
    	body { font-family: Arial, sans-serif; margin: 40px; }
        .header { font-size: 24px; color: #333; margin-bottom: 20px; }
        .file-list { list-style: none; padding: 0; }
        .file-item { padding: 8px; border-bottom: 1px solid #eee; }
        .file-item:hover { background: #f9f9f9; }
        .file-link { text-decoration: none; color: #0366d6; }
        .file-link:hover { text-decoration: underline; }
        .folder::before { content: "📁 "; }
        .file::before { content: "📄 "; }
    </style>
</head>
<body>
	<div class="header">Index of {{.Dir}}</div>
	<ul class="file-list">
		{{range .Files}}
		<li class="file-item">
			<a href="{{.Path}}" class="file-link {{if .IsDir}}folder{{else}}file{{end}}">{{.Name}}</a>
		</li>
		{{end}}
	</ul>
</body>
</html>
`

// FileInfo contains basic info about a file.
type FileInfo struct {
	Name  string // file name
	Path  string // file path
	IsDir bool   // is file a directory
}

// This is a re-implementation of default http.FileServer handler, written to render custom HTML/CSS.
func FileServer(root http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := root.Open(r.URL.Path)
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
				Path:  filepath.ToSlash(filepath.Join(r.URL.Path, file.Name())),
				IsDir: file.IsDir(),
			})
		}

		// rendering page
		tmpl := template.Must(template.New("dir").Parse(dirTemplate))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, struct {
			Dir   string
			Files []FileInfo
		}{
			Dir:   r.URL.Path,
			Files: fileList,
		})
	})
}
