package handlers

import (
	"context"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	StatusCode() int
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *responseWriter) Write(bytes []byte) (int, error) {
	return w.ResponseWriter.Write(bytes)
}

func (w *responseWriter) StatusCode() int {
	return w.statusCode
}

type AppContext string

var StatusCodeContext = AppContext("statusCode")

func WriteHeaderAndContext(w http.ResponseWriter, r *http.Request, statusCode int) {
	ctx := context.WithValue(r.Context(), StatusCodeContext, statusCode)
	*r = *(r.WithContext(ctx))

	w.WriteHeader(statusCode)
}
