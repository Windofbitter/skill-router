package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		enabledDir:  filepath.Join(baseDir, "commands"),
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
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		fm, _ := parser.ParseFrontmatter(string(content))
		name := fm.Name
		if name == "" {
			name = strings.TrimSuffix(entry.Name(), ".md")
		}

		skills = append(skills, model.Skill{
			Name:        name,
			Description: fm.Description,
			FileName:    entry.Name(),
			FilePath:    filePath,
			Enabled:     enabled,
		})
	}

	return skills, nil
}

func (s *SkillService) DisableSkill(fileName string) error {
	src := filepath.Join(s.enabledDir, fileName)
	dst := filepath.Join(s.disabledDir, fileName)

	if err := os.MkdirAll(s.disabledDir, 0755); err != nil {
		return err
	}

	return os.Rename(src, dst)
}

func (s *SkillService) EnableSkill(fileName string) error {
	src := filepath.Join(s.disabledDir, fileName)
	dst := filepath.Join(s.enabledDir, fileName)

	if err := os.MkdirAll(s.enabledDir, 0755); err != nil {
		return err
	}

	return os.Rename(src, dst)
}

func (s *SkillService) DeleteSkill(fileName string, enabled bool) error {
	var filePath string
	if enabled {
		filePath = filepath.Join(s.enabledDir, fileName)
	} else {
		filePath = filepath.Join(s.disabledDir, fileName)
	}
	return os.Remove(filePath)
}

func (s *SkillService) SaveSkill(fileName string, content []byte, overwrite bool) error {
	filePath := filepath.Join(s.enabledDir, fileName)

	if !overwrite {
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("file already exists")
		}
	}

	if err := os.MkdirAll(s.enabledDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, content, 0644)
}
