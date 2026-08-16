package webserve

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

func registerEmbeddedStaticRoutes(mux *http.ServeMux, content fs.FS) {
	// Astro Related Path
	mux.HandleFunc("GET /_astro/{filepath...}", func(w http.ResponseWriter, r *http.Request) {
		serveEmbeddedFile(w, r, content, "_astro", r.PathValue("filepath"))
	})

	// Generic web application dist paths
	mux.HandleFunc("GET /assets/{filepath}", func(w http.ResponseWriter, r *http.Request) {
		serveEmbeddedFile(w, r, content, "assets", r.PathValue("filepath"))
	})
	mux.HandleFunc("GET /favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		serveEmbeddedFile(w, r, content, "", "favicon.svg")
	})
	mux.HandleFunc("GET /apple-touch-icon.png", func(w http.ResponseWriter, r *http.Request) {
		serveEmbeddedFile(w, r, content, "", "apple-touch-icon.png")
	})
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		serveEmbeddedFile(w, r, content, "", "assets/robots.txt")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			http.NotFound(w, r)
			return
		}
		serveEmbeddedFile(w, r, content, "", "index.html")
	})
}

func serveEmbeddedFile(w http.ResponseWriter, r *http.Request, content fs.FS, prefix string, filepath string) {
	// If the path is empty or ends with "/", serve index.html
	if filepath == "" || strings.HasSuffix(filepath, "/") {
		filepath = path.Join(filepath, "index.html")
	}

	fullPath := path.Join(prefix, filepath)
	f, err := content.Open(fullPath)
	if err != nil {
		// If file not found, serve 404
		fmt.Printf("File not found: %s\n", fullPath)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, "read error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("%s: %d bytes\n", fullPath, len(data))

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f.(io.ReadSeeker))
}
