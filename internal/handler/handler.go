package handler

import (
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	service Service
}

func New(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Route() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /", h.HomePage)
	router.HandleFunc("GET /order", h.GetOrder)

	return router
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/page.html")
	if err != nil {
		http.Error(w, "Sorry, page unavailable", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Sorry, page rendering failed", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if content, ok := h.service.GetOrderByID(id); !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}
}
