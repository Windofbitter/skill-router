package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type SkillOverrides struct {
	Disabled []string `json:"disabled"`
}

var (
	overridesPath string
	overridesMu   sync.RWMutex
)

func Init(baseDir string) {
	overridesPath = filepath.Join(baseDir, "skill-overrides.json")
}

func LoadOverrides() (*SkillOverrides, error) {
	overridesMu.RLock()
	defer overridesMu.RUnlock()

	data, err := os.ReadFile(overridesPath)
	if os.IsNotExist(err) {
		return &SkillOverrides{Disabled: []string{}}, nil
	}
	if err != nil {
		return nil, err
	}

	var overrides SkillOverrides
	if err := json.Unmarshal(data, &overrides); err != nil {
		return nil, err
	}

	return &overrides, nil
}

func SaveOverrides(overrides *SkillOverrides) error {
	overridesMu.Lock()
	defer overridesMu.Unlock()

	data, err := json.MarshalIndent(overrides, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(overridesPath, data, 0644)
}

func IsPluginSkillDisabled(pluginName, skillName string) bool {
	overrides, err := LoadOverrides()
	if err != nil {
		return false
	}

	key := pluginName + ":" + skillName
	for _, disabled := range overrides.Disabled {
		if disabled == key {
			return true
		}
	}
	return false
}

func DisablePluginSkill(pluginName, skillName string) error {
	overrides, err := LoadOverrides()
	if err != nil {
		return err
	}

	key := pluginName + ":" + skillName

	// Check if already disabled
	for _, disabled := range overrides.Disabled {
		if disabled == key {
			return nil
		}
	}

	overrides.Disabled = append(overrides.Disabled, key)
	return SaveOverrides(overrides)
}

func EnablePluginSkill(pluginName, skillName string) error {
	overrides, err := LoadOverrides()
	if err != nil {
		return err
	}

	key := pluginName + ":" + skillName

	// Remove from disabled list
	newDisabled := make([]string, 0, len(overrides.Disabled))
	for _, disabled := range overrides.Disabled {
		if disabled != key {
			newDisabled = append(newDisabled, disabled)
		}
	}

	overrides.Disabled = newDisabled
	return SaveOverrides(overrides)
}
