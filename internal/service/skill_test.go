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

func TestDisableSkill(t *testing.T) {
	tmpDir := t.TempDir()
	enabledDir := filepath.Join(tmpDir, "commands")
	disabledDir := filepath.Join(tmpDir, "skills-disabled")
	os.MkdirAll(enabledDir, 0755)
	os.MkdirAll(disabledDir, 0755)

	skillContent := `---
name: my-skill
description: Test
---
Content`
	os.WriteFile(filepath.Join(enabledDir, "my-skill.md"), []byte(skillContent), 0644)

	svc := NewSkillService(tmpDir)
	err := svc.DisableSkill("my-skill.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file moved
	if _, err := os.Stat(filepath.Join(enabledDir, "my-skill.md")); !os.IsNotExist(err) {
		t.Error("file should be removed from enabled dir")
	}
	if _, err := os.Stat(filepath.Join(disabledDir, "my-skill.md")); err != nil {
		t.Error("file should exist in disabled dir")
	}
}

func TestEnableSkill(t *testing.T) {
	tmpDir := t.TempDir()
	enabledDir := filepath.Join(tmpDir, "commands")
	disabledDir := filepath.Join(tmpDir, "skills-disabled")
	os.MkdirAll(enabledDir, 0755)
	os.MkdirAll(disabledDir, 0755)

	skillContent := `---
name: my-skill
description: Test
---
Content`
	os.WriteFile(filepath.Join(disabledDir, "my-skill.md"), []byte(skillContent), 0644)

	svc := NewSkillService(tmpDir)
	err := svc.EnableSkill("my-skill.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file moved
	if _, err := os.Stat(filepath.Join(disabledDir, "my-skill.md")); !os.IsNotExist(err) {
		t.Error("file should be removed from disabled dir")
	}
	if _, err := os.Stat(filepath.Join(enabledDir, "my-skill.md")); err != nil {
		t.Error("file should exist in enabled dir")
	}
}

func TestDeleteSkill(t *testing.T) {
	tmpDir := t.TempDir()
	enabledDir := filepath.Join(tmpDir, "commands")
	os.MkdirAll(enabledDir, 0755)

	skillContent := `---
name: my-skill
description: Test
---
Content`
	os.WriteFile(filepath.Join(enabledDir, "my-skill.md"), []byte(skillContent), 0644)

	svc := NewSkillService(tmpDir)
	err := svc.DeleteSkill("my-skill.md", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(enabledDir, "my-skill.md")); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}
