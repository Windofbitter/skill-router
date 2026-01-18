package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListSkills(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	enabledDir := filepath.Join(tmpDir, "commands")
	disabledDir := filepath.Join(tmpDir, "skills-disabled")
	os.MkdirAll(enabledDir, 0755)
	os.MkdirAll(disabledDir, 0755)

	// Create test skill files
	skill1 := `---
name: skill-one
description: First skill
---
Content`
	os.WriteFile(filepath.Join(enabledDir, "skill-one.md"), []byte(skill1), 0644)

	skill2 := `---
name: skill-two
description: Second skill (disabled)
---
Content`
	os.WriteFile(filepath.Join(disabledDir, "skill-two.md"), []byte(skill2), 0644)

	svc := NewSkillService(tmpDir)
	skills, err := svc.ListSkills()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}

	// Check enabled skill
	var enabled, disabled bool
	for _, s := range skills {
		if s.Name == "skill-one" && s.Enabled {
			enabled = true
		}
		if s.Name == "skill-two" && !s.Enabled {
			disabled = true
		}
	}

	if !enabled {
		t.Error("skill-one should be enabled")
	}
	if !disabled {
		t.Error("skill-two should be disabled")
	}
}
