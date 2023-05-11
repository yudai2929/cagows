package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/sivchari/cagows/answer/model"
	"github.com/sivchari/cagows/answer/repository"
)

type Handler interface {
	List() http.HandlerFunc
	Get() http.HandlerFunc
	Add() http.HandlerFunc
	Complete() http.HandlerFunc
}

type handler struct {
	repo repository.Repository
}

func NewHandler(repo repository.Repository) Handler {
	return &handler{
		repo: repo,
	}
}

func (h *handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos := h.repo.List()
		if err := ResponseJSON(w, todos); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := GetID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todo, err := h.repo.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := ResponseJSON(w, todo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *handler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo model.Todo
		if err := DecodeJSON(r, &todo); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newtodo := h.repo.Add(&todo)
		w.WriteHeader(http.StatusCreated)
		if err := ResponseJSON(w, newtodo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *handler) Complete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Receive struct {
			ID int `json:"id"`
		}
		var receive Receive
		if err := DecodeJSON(r, &receive); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := h.repo.Complete(receive.ID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetID(r *http.Request) (int, error) {
	url := r.URL.Query()
	id := url.Get("id")
	if id == "" {
		return 0, errors.New("id is empty")
	}
	return strconv.Atoi(id)
}

func DecodeJSON(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.Unmarshal(body, v)
}

func ResponseJSON(w http.ResponseWriter, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(body); err != nil {
		return err
	}
	return nil
}
