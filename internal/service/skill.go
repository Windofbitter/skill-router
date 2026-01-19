package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wind/skill-router/internal/config"
	"github.com/wind/skill-router/internal/model"
	"github.com/wind/skill-router/internal/parser"
)

type SkillService struct {
	baseDir     string
	enabledDir  string
	disabledDir string
	pluginsDir  string
}

func NewSkillService(baseDir string) *SkillService {
	return &SkillService{
		baseDir:     baseDir,
		enabledDir:  filepath.Join(baseDir, "skills"),
		disabledDir: filepath.Join(baseDir, "skills-disabled"),
		pluginsDir:  filepath.Join(baseDir, "plugins", "cache"),
	}
}

func (s *SkillService) ListSkills() ([]model.Skill, error) {
	var skills []model.Skill

	// Scan user enabled skills
	enabledSkills, err := s.scanUserDir(s.enabledDir, true)
	if err != nil {
		return nil, err
	}
	skills = append(skills, enabledSkills...)

	// Scan user disabled skills
	disabledSkills, err := s.scanUserDir(s.disabledDir, false)
	if err != nil {
		return nil, err
	}
	skills = append(skills, disabledSkills...)

	// Scan plugin skills
	pluginSkills, err := s.scanPlugins()
	if err != nil {
		return nil, err
	}
	skills = append(skills, pluginSkills...)

	return skills, nil
}

func (s *SkillService) scanUserDir(dir string, enabled bool) ([]model.Skill, error) {
	var skills []model.Skill

	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return skills, nil
	}
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skill := s.readSkillDir(filepath.Join(dir, entry.Name()), entry.Name())
		if skill != nil {
			skill.Enabled = enabled
			skill.Source = "user"
			skills = append(skills, *skill)
		}
	}

	return skills, nil
}

func (s *SkillService) scanPlugins() ([]model.Skill, error) {
	var skills []model.Skill

	// Structure: plugins/cache/<org>/<plugin>/<version>/skills/<skill-name>/SKILL.md
	orgs, err := os.ReadDir(s.pluginsDir)
	if os.IsNotExist(err) {
		return skills, nil
	}
	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		if !org.IsDir() {
			continue
		}

		orgPath := filepath.Join(s.pluginsDir, org.Name())
		plugins, err := os.ReadDir(orgPath)
		if err != nil {
			continue
		}

		for _, plugin := range plugins {
			if !plugin.IsDir() {
				continue
			}

			pluginPath := filepath.Join(orgPath, plugin.Name())
			versions, err := os.ReadDir(pluginPath)
			if err != nil {
				continue
			}

			// Use the latest version (last in sorted order)
			var latestVersion string
			for _, v := range versions {
				if v.IsDir() {
					latestVersion = v.Name()
				}
			}

			if latestVersion == "" {
				continue
			}

			skillsPath := filepath.Join(pluginPath, latestVersion, "skills")
			skillDirs, err := os.ReadDir(skillsPath)
			if err != nil {
				continue
			}

			pluginName := plugin.Name()

			for _, skillDir := range skillDirs {
				if !skillDir.IsDir() {
					continue
				}

				skill := s.readSkillDir(filepath.Join(skillsPath, skillDir.Name()), skillDir.Name())
				if skill != nil {
					skill.Enabled = !config.IsPluginSkillDisabled(pluginName, skillDir.Name())
					skill.Source = "plugin"
					skill.PluginName = pluginName
					skills = append(skills, *skill)
				}
			}
		}
	}

	return skills, nil
}

func (s *SkillService) readSkillDir(skillDir, dirName string) *model.Skill {
	skillFile := filepath.Join(skillDir, "SKILL.md")

	content, err := os.ReadFile(skillFile)
	if err != nil {
		// Try lowercase
		skillFile = filepath.Join(skillDir, "skill.md")
		content, err = os.ReadFile(skillFile)
		if err != nil {
			return nil
		}
	}

	fm, _ := parser.ParseFrontmatter(string(content))
	name := fm.Name
	if name == "" {
		name = dirName
	}

	return &model.Skill{
		Name:        name,
		Description: fm.Description,
		FileName:    dirName,
		FilePath:    skillDir,
	}
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

func (s *SkillService) SaveSkill(skillDirName string, content []byte, overwrite bool) error {
	skillDir := filepath.Join(s.enabledDir, skillDirName)
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

func (s *SkillService) DeletePlugin(pluginName string) error {
	// Find and delete the plugin directory
	orgs, err := os.ReadDir(s.pluginsDir)
	if err != nil {
		return err
	}

	for _, org := range orgs {
		if !org.IsDir() {
			continue
		}

		pluginPath := filepath.Join(s.pluginsDir, org.Name(), pluginName)
		if _, err := os.Stat(pluginPath); err == nil {
			return os.RemoveAll(pluginPath)
		}
	}

	return fmt.Errorf("plugin not found: %s", pluginName)
}
