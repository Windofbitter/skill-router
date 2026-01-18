package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type File struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
}

var repoURLRegex = regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`)

func FetchSkillFiles(repoURL string) ([]File, error) {
	matches := repoURLRegex.FindStringSubmatch(repoURL)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid github URL")
	}

	owner, repo := matches[1], matches[2]
	repo = strings.TrimSuffix(repo, ".git")

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/.claude/commands", owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("no .claude/commands directory found")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github API error: %d", resp.StatusCode)
	}

	var files []File
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	var mdFiles []File
	for _, f := range files {
		if f.Type == "file" && strings.HasSuffix(f.Name, ".md") {
			mdFiles = append(mdFiles, f)
		}
	}

	return mdFiles, nil
}

func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
