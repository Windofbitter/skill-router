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

var (
	repoURLRegex            = regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`) // owner/repo
	defaultGitHubAPIBaseURL = "https://api.github.com"
)

func FetchSkillFiles(repoURL string) ([]File, error) {
	return fetchSkillFiles(repoURL, defaultGitHubAPIBaseURL, http.DefaultClient)
}

func FetchSkillDirs(repoURL string) (basePath string, dirs []File, err error) {
	return fetchSkillDirs(repoURL, defaultGitHubAPIBaseURL, http.DefaultClient)
}

func FetchSkillFile(repoURL, basePath, skillDir string) (File, error) {
	return fetchSkillFile(repoURL, basePath, skillDir, defaultGitHubAPIBaseURL, http.DefaultClient)
}

func fetchSkillFiles(repoURL, apiBaseURL string, client *http.Client) ([]File, error) {
	basePath, dirs, err := fetchSkillDirs(repoURL, apiBaseURL, client)
	if err != nil {
		return nil, err
	}

	var skillFiles []File
	for _, d := range dirs {
		f, err := fetchSkillFile(repoURL, basePath, d.Name, apiBaseURL, client)
		if err != nil {
			return nil, err
		}

		skillFiles = append(skillFiles, File{
			Name:        d.Name,
			DownloadURL: f.DownloadURL,
			Type:        "file",
		})
	}

	return skillFiles, nil
}

func fetchSkillDirs(repoURL, apiBaseURL string, client *http.Client) (basePath string, dirs []File, err error) {
	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return "", nil, err
	}

	tryPaths := []string{".claude/skills", "skills"}

	for i, p := range tryPaths {
		apiURL := fmt.Sprintf("%s/repos/%s/%s/contents/%s", strings.TrimSuffix(apiBaseURL, "/"), owner, repo, p)

		var entries []File
		status, err := fetchJSON(client, apiURL, &entries)
		if err != nil {
			return "", nil, err
		}

		if status == http.StatusNotFound {
			if i == len(tryPaths)-1 {
				return "", nil, fmt.Errorf("no skills directory found")
			}
			continue
		}

		if status != http.StatusOK {
			return "", nil, fmt.Errorf("github API error: %d", status)
		}

		var dirEntries []File
		for _, e := range entries {
			if e.Type == "dir" {
				dirEntries = append(dirEntries, e)
			}
		}

		return p, dirEntries, nil
	}

	return "", nil, fmt.Errorf("no skills directory found")
}

func fetchSkillFile(repoURL, basePath, skillDir, apiBaseURL string, client *http.Client) (File, error) {
	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return File{}, err
	}

	tryFiles := []string{"SKILL.md", "skill.md"}
	for i, name := range tryFiles {
		path := fmt.Sprintf("%s/%s/%s", basePath, skillDir, name)
		apiURL := fmt.Sprintf("%s/repos/%s/%s/contents/%s", strings.TrimSuffix(apiBaseURL, "/"), owner, repo, path)

		var file File
		status, err := fetchJSON(client, apiURL, &file)
		if err != nil {
			return File{}, err
		}

		if status == http.StatusNotFound {
			if i == len(tryFiles)-1 {
				return File{}, fmt.Errorf("no skill file found for %s", skillDir)
			}
			continue
		}

		if status != http.StatusOK {
			return File{}, fmt.Errorf("github API error: %d", status)
		}

		return file, nil
	}

	return File{}, fmt.Errorf("no skill file found for %s", skillDir)
}

func parseRepoURL(repoURL string) (owner, repo string, err error) {
	matches := repoURLRegex.FindStringSubmatch(repoURL)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("invalid github URL")
	}

	owner, repo = matches[1], matches[2]
	repo = strings.TrimSuffix(repo, ".git")
	return owner, repo, nil
}

func fetchJSON(client *http.Client, url string, out any) (statusCode int, err error) {
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode
	if statusCode != http.StatusOK {
		return statusCode, nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return 0, err
	}

	return statusCode, nil
}

func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
