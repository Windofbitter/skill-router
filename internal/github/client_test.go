package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchSkillFiles_PrefersDotClaudeSkills(t *testing.T) {
	var hitsDotClaudeSkills int
	var hitsSkills int

	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/acme/myrepo/contents/.claude/skills":
			hitsDotClaudeSkills++
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]File{{Name: "foo", Type: "dir"}})
			return
		case "/repos/acme/myrepo/contents/skills":
			hitsSkills++
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]File{{Name: "bar", Type: "dir"}})
			return
		case "/repos/acme/myrepo/contents/.claude/skills/foo/SKILL.md":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(File{Name: "SKILL.md", Type: "file", DownloadURL: ts.URL + "/download/foo"})
			return
		case "/download/foo":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("# foo\n"))
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer ts.Close()

	files, err := fetchSkillFiles("https://github.com/acme/myrepo", ts.URL, ts.Client())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hitsDotClaudeSkills != 1 {
		t.Fatalf("expected .claude/skills to be queried once, got %d", hitsDotClaudeSkills)
	}
	if hitsSkills != 0 {
		t.Fatalf("expected skills/ not to be queried, got %d", hitsSkills)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 skill file, got %d", len(files))
	}
	if files[0].Name != "foo" {
		t.Fatalf("expected skill name foo, got %q", files[0].Name)
	}
	if files[0].DownloadURL != ts.URL+"/download/foo" {
		t.Fatalf("expected download url %q, got %q", ts.URL+"/download/foo", files[0].DownloadURL)
	}

	body, err := DownloadFile(files[0].DownloadURL)
	if err != nil {
		t.Fatalf("expected download to succeed, got %v", err)
	}
	if string(body) != "# foo\n" {
		t.Fatalf("unexpected download content: %q", string(body))
	}
}

func TestFetchSkillFiles_FallsBackToSkills(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/acme/myrepo/contents/.claude/skills":
			http.NotFound(w, r)
			return
		case "/repos/acme/myrepo/contents/skills":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]File{{Name: "foo", Type: "dir"}})
			return
		case "/repos/acme/myrepo/contents/skills/foo/SKILL.md":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(File{Name: "SKILL.md", Type: "file", DownloadURL: ts.URL + "/download/foo"})
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer ts.Close()

	files, err := fetchSkillFiles("https://github.com/acme/myrepo", ts.URL, ts.Client())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 skill file, got %d", len(files))
	}
	if files[0].Name != "foo" {
		t.Fatalf("expected skill name foo, got %q", files[0].Name)
	}
	if files[0].DownloadURL != ts.URL+"/download/foo" {
		t.Fatalf("expected download url %q, got %q", ts.URL+"/download/foo", files[0].DownloadURL)
	}
}

func TestFetchSkillFiles_UsesSkillMdFallback(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/acme/myrepo/contents/.claude/skills":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]File{{Name: "foo", Type: "dir"}})
			return
		case "/repos/acme/myrepo/contents/.claude/skills/foo/SKILL.md":
			http.NotFound(w, r)
			return
		case "/repos/acme/myrepo/contents/.claude/skills/foo/skill.md":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(File{Name: "skill.md", Type: "file", DownloadURL: ts.URL + "/download/foo"})
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer ts.Close()

	files, err := fetchSkillFiles("https://github.com/acme/myrepo", ts.URL, ts.Client())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 skill file, got %d", len(files))
	}
	if files[0].Name != "foo" {
		t.Fatalf("expected skill name foo, got %q", files[0].Name)
	}
	if files[0].DownloadURL != ts.URL+"/download/foo" {
		t.Fatalf("expected download url %q, got %q", ts.URL+"/download/foo", files[0].DownloadURL)
	}
}
