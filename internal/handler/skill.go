package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/wind/skill-router/internal/github"
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

func (h *SkillHandler) Upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	overwrite := r.FormValue("overwrite") == "true"

	if err := h.svc.SaveSkill(header.Filename, content, overwrite); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type InstallRequest struct {
	URL string `json:"url"`
}

func (h *SkillHandler) Install(w http.ResponseWriter, r *http.Request) {
	var req InstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	files, err := github.FetchSkillFiles(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	installed := 0
	for _, f := range files {
		content, err := github.DownloadFile(f.DownloadURL)
		if err != nil {
			continue
		}
		if err := h.svc.SaveSkill(f.Name, content, false); err != nil {
			continue
		}
		installed++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"installed": installed})
}
