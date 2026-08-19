package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var filePath = flag.String("d", "", "Path to file directory")
var uploadDir = filePath

func handleIndex(fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, `<form method="POST" enctype="multipart/form-data">
	<input type="file" name="file">
	<input type="submit" value="Upload">
</form>`)
		}
		fileServer.ServeHTTP(w, r)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<form method="POST" enctype="multipart/form-data">
	<input type="file" name="file">
	<input type="submit" value="Upload">
</form>`)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	dstPath := filepath.Join(*uploadDir, filepath.Base(header.Filename))
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Uploaded file: " + header.Filename))
}

func main() {
	serverPort := flag.Int("p", 5000, "Server port")
	flag.Parse()

	fileServer := http.FileServer(http.Dir("./" + *filePath))

	http.HandleFunc("/", handleIndex(fileServer))
	http.HandleFunc("/upload", handleUpload)

	log.Printf("Listening on http://localhost:%d\n", *serverPort)
	http.ListenAndServe(":"+strconv.Itoa(*serverPort), nil)
}
