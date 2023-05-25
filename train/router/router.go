package router

import (
	"github.com/go-chi/chi"

	"github.com/sivchari/cagows/train/handler"
	"github.com/sivchari/cagows/train/repository"
)

func Routing(repo repository.Repository) *chi.Mux {
	h := handler.NewHandler(repo)
	r := chi.NewRouter()
	r.Get("/list", h.List())
	r.Get("/get", h.Get())
	// TODO: /add と /complete を追加(POST method)
	r.Post("/add", h.Add())
	r.Post("/complete", h.Complete())
	return r
}
