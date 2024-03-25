package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/NurfitraPujo/image-processor/internal/handlers"
	"github.com/NurfitraPujo/image-processor/internal/images"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	sloghttp "github.com/samber/slog-http"
)

func StartServer() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := new(handlers.MiddlewareMux)
	handler := sloghttp.Recovery(mux)
	handler = sloghttp.New(logger)(handler)

	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.Handle("GET /images/", http.StripPrefix("/images/", http.FileServer(http.Dir("public/images"))))
	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		aliasFileName := r.FormValue("alias")
		uploadedFile, handler, err := r.FormFile("file")
		if err != nil {
			slog.Error("error when parsing form data", slog.String("error", err.Error()))

			http.Error(w, "unexpected error when parsing form data", http.StatusInternalServerError)
			return
		}
		defer uploadedFile.Close()

		err = images.SaveImage(aliasFileName, *handler, uploadedFile)
		if err != nil {
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
