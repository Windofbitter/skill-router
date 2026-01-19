package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/wind/skill-router/internal/config"
	"github.com/wind/skill-router/internal/github"
	"github.com/wind/skill-router/internal/parser"
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

	fm, _ := parser.ParseFrontmatter(string(content))
	skillDir := strings.TrimSpace(fm.Name)
	if skillDir == "" {
		base := filepath.Base(header.Filename)
		skillDir = strings.TrimSuffix(base, filepath.Ext(base))
	}

	if err := h.svc.SaveSkill(skillDir, content, overwrite); err != nil {
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

// Plugin skill handlers - these modify the override config file

func (h *SkillHandler) DisablePluginSkill(w http.ResponseWriter, r *http.Request) {
	// URL format: /api/plugins/{pluginName}/skills/{skillName}/disable
	path := strings.TrimPrefix(r.URL.Path, "/api/plugins/")
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	pluginName := parts[0]
	skillName := parts[2]

	if err := config.DisablePluginSkill(pluginName, skillName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SkillHandler) EnablePluginSkill(w http.ResponseWriter, r *http.Request) {
	// URL format: /api/plugins/{pluginName}/skills/{skillName}/enable
	path := strings.TrimPrefix(r.URL.Path, "/api/plugins/")
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	pluginName := parts[0]
	skillName := parts[2]

	if err := config.EnablePluginSkill(pluginName, skillName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SkillHandler) DisablePlugin(w http.ResponseWriter, r *http.Request) {
	// URL format: /api/plugins/{pluginName}/disable
	path := strings.TrimPrefix(r.URL.Path, "/api/plugins/")
	pluginName := strings.TrimSuffix(path, "/disable")

	if err := config.DisablePlugin(pluginName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SkillHandler) EnablePlugin(w http.ResponseWriter, r *http.Request) {
	// URL format: /api/plugins/{pluginName}/enable
	path := strings.TrimPrefix(r.URL.Path, "/api/plugins/")
	pluginName := strings.TrimSuffix(path, "/enable")

	if err := config.EnablePlugin(pluginName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SkillHandler) DeletePlugin(w http.ResponseWriter, r *http.Request) {
	// URL format: /api/plugins/{pluginName}
	pluginName := strings.TrimPrefix(r.URL.Path, "/api/plugins/")

	if err := h.svc.DeletePlugin(pluginName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
