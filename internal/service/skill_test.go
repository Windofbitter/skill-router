package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListSkills(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	enabledDir := filepath.Join(tmpDir, "skills")
	disabledDir := filepath.Join(tmpDir, "skills-disabled")

	// Create skill directories with SKILL.md files
	skill1Dir := filepath.Join(enabledDir, "skill-one")
	os.MkdirAll(skill1Dir, 0755)
	skill1 := `---
name: skill-one
description: First skill
---
Content`
	os.WriteFile(filepath.Join(skill1Dir, "SKILL.md"), []byte(skill1), 0644)

	skill2Dir := filepath.Join(disabledDir, "skill-two")
	os.MkdirAll(skill2Dir, 0755)
	skill2 := `---
name: skill-two
description: Second skill (disabled)
---
Content`
	os.WriteFile(filepath.Join(skill2Dir, "SKILL.md"), []byte(skill2), 0644)

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
	enabledDir := filepath.Join(tmpDir, "skills")
	disabledDir := filepath.Join(tmpDir, "skills-disabled")

	// Create skill directory
	skillDir := filepath.Join(enabledDir, "my-skill")
	os.MkdirAll(skillDir, 0755)
	skillContent := `---
name: my-skill
description: Test
---
Content`
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0644)

	svc := NewSkillService(tmpDir)
	err := svc.DisableSkill("my-skill")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify directory moved
	if _, err := os.Stat(filepath.Join(enabledDir, "my-skill")); !os.IsNotExist(err) {
		t.Error("skill should be removed from enabled dir")
	}
	if _, err := os.Stat(filepath.Join(disabledDir, "my-skill", "SKILL.md")); err != nil {
		t.Error("skill should exist in disabled dir")
	}
}

func TestEnableSkill(t *testing.T) {
	tmpDir := t.TempDir()
	enabledDir := filepath.Join(tmpDir, "skills")
	disabledDir := filepath.Join(tmpDir, "skills-disabled")

	// Create disabled skill directory
	skillDir := filepath.Join(disabledDir, "my-skill")
	os.MkdirAll(skillDir, 0755)
	skillContent := `---
name: my-skill
description: Test
---
Content`
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0644)

	svc := NewSkillService(tmpDir)
	err := svc.EnableSkill("my-skill")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify directory moved
	if _, err := os.Stat(filepath.Join(disabledDir, "my-skill")); !os.IsNotExist(err) {
		t.Error("skill should be removed from disabled dir")
	}
	if _, err := os.Stat(filepath.Join(enabledDir, "my-skill", "SKILL.md")); err != nil {
		t.Error("skill should exist in enabled dir")
	}
}

func TestDeleteSkill(t *testing.T) {
	tmpDir := t.TempDir()
	enabledDir := filepath.Join(tmpDir, "skills")

	// Create skill directory
	skillDir := filepath.Join(enabledDir, "my-skill")
	os.MkdirAll(skillDir, 0755)
	skillContent := `---
name: my-skill
description: Test
---
Content`
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0644)

	svc := NewSkillService(tmpDir)
	err := svc.DeleteSkill("my-skill", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(enabledDir, "my-skill")); !os.IsNotExist(err) {
		t.Error("skill directory should be deleted")
	}
}

func TestSaveSkill_WritesSkillBundle(t *testing.T) {
	tmpDir := t.TempDir()

	content := []byte(`---
name: my-skill
description: Test
---
Content`)

	svc := NewSkillService(tmpDir)
	if err := svc.SaveSkill("my-skill", content, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(tmpDir, "skills", "my-skill", "SKILL.md"))
	if err != nil {
		t.Fatalf("expected skill to be saved under skills/my-skill/SKILL.md: %v", err)
	}
	if string(got) != string(content) {
		t.Fatalf("unexpected saved content: %q", string(got))
	}
}
