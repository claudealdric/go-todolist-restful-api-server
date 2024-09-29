package api

import (
	"fmt"
	"net/http"
)

type Router struct {
	http.ServeMux
}

func NewRouter(s *Server) *Router {
	r := Router{}
	r.Get("/{$}", s.HandleRoot)
	r.Get("/tasks", s.HandleGetTasks)
	r.Get("/tasks/{id}", s.HandleGetTaskById)
	r.Patch("/tasks/{id}", s.HandlePatchTask)
	r.Post("/tasks", s.HandlePostTask)
	r.Post("/users", s.HandlePostUser)
	r.Delete("/tasks/{id}", s.HandleDeleteTask)
	return &r
}

func (r *Router) Delete(pattern string, handlerFunc http.HandlerFunc) {
	r.getHandlerFuncPattern(http.MethodDelete, pattern, handlerFunc)
}

func (r *Router) Get(pattern string, handlerFunc http.HandlerFunc) {
	r.getHandlerFuncPattern(http.MethodGet, pattern, handlerFunc)
}

func (r *Router) Patch(pattern string, handlerFunc http.HandlerFunc) {
	r.getHandlerFuncPattern(http.MethodPatch, pattern, handlerFunc)
}

func (r *Router) Post(pattern string, handlerFunc http.HandlerFunc) {
	r.getHandlerFuncPattern(http.MethodPost, pattern, handlerFunc)
}

func (r *Router) getHandlerFuncPattern(
	method, pattern string,
	handlerFunc http.HandlerFunc,
) {
	r.HandleFunc(fmt.Sprintf("%s %s", method, pattern), handlerFunc)
}
