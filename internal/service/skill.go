package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wind/skill-router/internal/model"
	"github.com/wind/skill-router/internal/parser"
)

type SkillService struct {
	baseDir     string
	enabledDir  string
	disabledDir string
}

func NewSkillService(baseDir string) *SkillService {
	return &SkillService{
		baseDir:     baseDir,
		enabledDir:  filepath.Join(baseDir, "skills"),
		disabledDir: filepath.Join(baseDir, "skills-disabled"),
	}
}

func (s *SkillService) ListSkills() ([]model.Skill, error) {
	var skills []model.Skill

	// Scan enabled skills
	enabledSkills, err := s.scanDir(s.enabledDir, true)
	if err != nil {
		return nil, err
	}
	skills = append(skills, enabledSkills...)

	// Scan disabled skills
	disabledSkills, err := s.scanDir(s.disabledDir, false)
	if err != nil {
		return nil, err
	}
	skills = append(skills, disabledSkills...)

	return skills, nil
}

func (s *SkillService) scanDir(dir string, enabled bool) ([]model.Skill, error) {
	var skills []model.Skill

	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return skills, nil
	}
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// Skills are directories containing SKILL.md
		if !entry.IsDir() {
			continue
		}

		skillDir := filepath.Join(dir, entry.Name())
		skillFile := filepath.Join(skillDir, "SKILL.md")

		content, err := os.ReadFile(skillFile)
		if err != nil {
			// Try lowercase skill.md as fallback
			skillFile = filepath.Join(skillDir, "skill.md")
			content, err = os.ReadFile(skillFile)
			if err != nil {
				continue
			}
		}

		fm, _ := parser.ParseFrontmatter(string(content))
		name := fm.Name
		if name == "" {
			name = entry.Name()
		}

		skills = append(skills, model.Skill{
			Name:        name,
			Description: fm.Description,
			FileName:    entry.Name(), // Directory name
			FilePath:    skillDir,
			Enabled:     enabled,
		})
	}

	return skills, nil
}

func (s *SkillService) DisableSkill(dirName string) error {
	src := filepath.Join(s.enabledDir, dirName)
	dst := filepath.Join(s.disabledDir, dirName)

	if err := os.MkdirAll(s.disabledDir, 0755); err != nil {
		return err
	}

	return os.Rename(src, dst)
}

func (s *SkillService) EnableSkill(dirName string) error {
	src := filepath.Join(s.disabledDir, dirName)
	dst := filepath.Join(s.enabledDir, dirName)

	if err := os.MkdirAll(s.enabledDir, 0755); err != nil {
		return err
	}

	return os.Rename(src, dst)
}

func (s *SkillService) DeleteSkill(dirName string, enabled bool) error {
	var dirPath string
	if enabled {
		dirPath = filepath.Join(s.enabledDir, dirName)
	} else {
		dirPath = filepath.Join(s.disabledDir, dirName)
	}
	return os.RemoveAll(dirPath)
}

func (s *SkillService) SaveSkill(fileName string, content []byte, overwrite bool) error {
	// Extract skill name from filename (remove .md extension)
	skillName := fileName
	if len(skillName) > 3 && skillName[len(skillName)-3:] == ".md" {
		skillName = skillName[:len(skillName)-3]
	}

	skillDir := filepath.Join(s.enabledDir, skillName)
	skillFile := filepath.Join(skillDir, "SKILL.md")

	if !overwrite {
		if _, err := os.Stat(skillDir); err == nil {
			return fmt.Errorf("skill already exists")
		}
	}

	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(skillFile, content, 0644)
}
