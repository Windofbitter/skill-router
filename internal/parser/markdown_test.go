package parser

import (
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	content := `---
name: test-skill
description: A test skill
---

# Content here
`
	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.Name != "test-skill" {
		t.Errorf("expected name 'test-skill', got '%s'", fm.Name)
	}
	if fm.Description != "A test skill" {
		t.Errorf("expected description 'A test skill', got '%s'", fm.Description)
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	content := `# Just markdown without frontmatter`
	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.Name != "" {
		t.Errorf("expected empty name, got '%s'", fm.Name)
	}
}
