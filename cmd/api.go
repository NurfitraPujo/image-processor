package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NurfitraPujo/image-processor/internal/handlers"
	sloghttp "github.com/samber/slog-http"
)

func StartServer() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := new(handlers.MiddlewareMux)
	handler := sloghttp.Recovery(mux)
	handler = sloghttp.New(logger)(handler)

	imgDir := http.FileServer(http.Dir("public/images"))
	mux.Handle("/images/", http.StripPrefix("/images/", imgDir))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

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

	logger.Info(fmt.Sprintf("service started in port %s", port))
	http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}
