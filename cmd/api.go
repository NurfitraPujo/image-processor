package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/NurfitraPujo/image-processor/internal/handlers"
	"github.com/NurfitraPujo/image-processor/internal/images"
	"github.com/NurfitraPujo/image-processor/internal/middlewares"
	"github.com/go-chi/chi/v5"
	chiMiddlewares "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	slogchi "github.com/samber/slog-chi"
)

func StartServer() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := chi.NewMux()

	mux.Use(chiMiddlewares.RealIP)
	mux.Use(slogchi.New(logger))
	mux.Use(middlewares.PrometheusHttpMiddleware)
	mux.Use(chiMiddlewares.Recoverer)

	handlers.FileServer(mux, "/images/", http.Dir("public/images"))

	mux.Get("/metrics", promhttp.Handler().ServeHTTP)
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
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
	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
