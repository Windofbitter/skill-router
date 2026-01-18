package parser

import (
	"regexp"
	"strings"
)

type Frontmatter struct {
	Name        string
	Description string
}

var frontmatterRegex = regexp.MustCompile(`(?s)^---\n(.+?)\n---`)

func ParseFrontmatter(content string) (Frontmatter, error) {
	var fm Frontmatter

	matches := frontmatterRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return fm, nil
	}

	lines := strings.Split(matches[1], "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "name":
			fm.Name = value
		case "description":
			fm.Description = value
		}
	}

	return fm, nil
}
