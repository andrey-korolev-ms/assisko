package handlers

import (
	"assisko/models"
	"assisko/services"
	"encoding/json"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type Handler struct {
	service *services.Service
	session *scs.SessionManager
}

func NewHandler(svc *services.Service, session *scs.SessionManager) *Handler {
	return &Handler{service: svc, session: session}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.session.Put(r.Context(), "userName", "Andrey")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged_in"})
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	// пока просто примем куку
	name := h.session.GetString(r.Context(), "userName")
	if name == "" {
		http.Error(w, "not a logged in", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"user": name})
}

func (h *Handler) SaveData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() //контекст пока возьмем из http запроса
	var data models.Human
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.service.SaveData(ctx, data); err != nil {
		http.Error(w, "faild to save", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

}
func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("GET /profile", h.Profile)
	mux.HandleFunc("POST /data", h.SaveData)
	return mux
}
