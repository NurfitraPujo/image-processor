package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func StartServer() {
	mux := http.NewServeMux()

	imgDir := http.FileServer(http.Dir("public/images"))
	mux.Handle("/images/", http.StripPrefix("/images/", imgDir))

	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		aliasFileName := r.FormValue("alias")
		uploadedFile, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer uploadedFile.Close()

		dir, err := os.Getwd()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filename := handler.Filename
		if aliasFileName != "" {
			filename = fmt.Sprintf("%s%s", aliasFileName, filepath.Ext(handler.Filename))
		}

		fileLocation := filepath.Join(dir, "public/images", filename)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, uploadedFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("uploaded"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
