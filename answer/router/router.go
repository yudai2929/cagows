package router

import (
	"github.com/go-chi/chi"

	"github.com/sivchari/cagows/answer/handler"
	"github.com/sivchari/cagows/answer/repository"
)

func Routing(repo repository.Repository) *chi.Mux {
	h := handler.NewHandler(repo)
	r := chi.NewRouter()
	r.Get("/list", h.List())
	r.Get("/get", h.Get())
	r.Post("/add", h.Add())
	r.Post("/complete", h.Complete())
	return r
}
