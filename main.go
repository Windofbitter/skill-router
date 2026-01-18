package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wind/skill-router/internal/handler"
	"github.com/wind/skill-router/internal/service"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	claudeDir := filepath.Join(homeDir, ".claude")

	svc := service.NewSkillService(claudeDir)
	h := handler.NewSkillHandler(svc)

	http.HandleFunc("/api/skills", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h.List(w, r)
		}
	})

	http.HandleFunc("/api/skills/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h.Upload(w, r)
		}
	})

	http.HandleFunc("/api/skills/install", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h.Install(w, r)
		}
	})

	http.HandleFunc("/api/skills/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/disable") && r.Method == "POST":
			h.Disable(w, r)
		case strings.HasSuffix(path, "/enable") && r.Method == "POST":
			h.Enable(w, r)
		case r.Method == "DELETE":
			h.Delete(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	addr := ":9527"
	url := "http://localhost" + addr

	fmt.Printf("Skill Router running at %s\n", url)

	// Open browser
	go openBrowser(url)

	http.ListenAndServe(addr, nil)
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}
	if cmd != nil {
		cmd.Run()
	}
}
