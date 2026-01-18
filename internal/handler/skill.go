package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wind/skill-router/internal/service"
)

type SkillHandler struct {
	svc *service.SkillService
}

func NewSkillHandler(svc *service.SkillService) *SkillHandler {
	return &SkillHandler{svc: svc}
}

func (h *SkillHandler) List(w http.ResponseWriter, r *http.Request) {
	skills, err := h.svc.ListSkills()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skills)
}

func (h *SkillHandler) Disable(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	fileName = strings.TrimSuffix(fileName, "/disable")

	if err := h.svc.DisableSkill(fileName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SkillHandler) Enable(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	fileName = strings.TrimSuffix(fileName, "/enable")

	if err := h.svc.EnableSkill(fileName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SkillHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	enabled := r.URL.Query().Get("enabled") == "true"

	if err := h.svc.DeleteSkill(fileName, enabled); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
