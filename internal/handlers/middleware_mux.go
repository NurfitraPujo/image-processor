package handlers

import (
	"container/list"
	"net/http"
)

type MiddlewareType func(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request))

type MiddlewareMux struct {
	http.ServeMux
	middlewares list.List
}

func (mux *MiddlewareMux) AppendMiddleware(middleware MiddlewareType) {
	mux.middlewares.PushBack(MiddlewareType(middleware))
}

func (mux *MiddlewareMux) PrependMiddleware(middleware MiddlewareType) {
	mux.middlewares.PushFront(MiddlewareType(middleware))
}

func (mux *MiddlewareMux) nextMiddleware(el *list.Element) func(w http.ResponseWriter, req *http.Request) {
	if el != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			el.Value.(MiddlewareType)(w, req, mux.nextMiddleware(el.Next()))
		}
	}
	return mux.ServeMux.ServeHTTP
}

func (mux *MiddlewareMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux.nextMiddleware(mux.middlewares.Front())(w, req)
}
