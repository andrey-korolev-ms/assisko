package handlers

import (
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
	json.NewEncoder(w).Enconde(map[string]string{"status": "logged_in"})
}
