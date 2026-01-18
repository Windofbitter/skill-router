package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SkillOverrides struct {
	Disabled        []string `json:"disabled"`
	DisabledPlugins []string `json:"disabledPlugins"`
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
		return &SkillOverrides{Disabled: []string{}, DisabledPlugins: []string{}}, nil
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

	// Check if entire plugin is disabled
	for _, disabled := range overrides.DisabledPlugins {
		if disabled == pluginName {
			return true
		}
	}

	// Check if individual skill is disabled
	key := pluginName + ":" + skillName
	for _, disabled := range overrides.Disabled {
		if disabled == key {
			return true
		}
	}
	return false
}

func IsPluginDisabled(pluginName string) bool {
	overrides, err := LoadOverrides()
	if err != nil {
		return false
	}

	for _, disabled := range overrides.DisabledPlugins {
		if disabled == pluginName {
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

func DisablePlugin(pluginName string) error {
	overrides, err := LoadOverrides()
	if err != nil {
		return err
	}

	// Check if already disabled
	for _, disabled := range overrides.DisabledPlugins {
		if disabled == pluginName {
			return nil
		}
	}

	overrides.DisabledPlugins = append(overrides.DisabledPlugins, pluginName)

	// Also remove individual skill overrides for this plugin (they're now redundant)
	newDisabled := make([]string, 0, len(overrides.Disabled))
	for _, disabled := range overrides.Disabled {
		if !strings.HasPrefix(disabled, pluginName+":") {
			newDisabled = append(newDisabled, disabled)
		}
	}
	overrides.Disabled = newDisabled

	return SaveOverrides(overrides)
}

func EnablePlugin(pluginName string) error {
	overrides, err := LoadOverrides()
	if err != nil {
		return err
	}

	// Remove from disabled plugins list
	newDisabledPlugins := make([]string, 0, len(overrides.DisabledPlugins))
	for _, disabled := range overrides.DisabledPlugins {
		if disabled != pluginName {
			newDisabledPlugins = append(newDisabledPlugins, disabled)
		}
	}

	overrides.DisabledPlugins = newDisabledPlugins
	return SaveOverrides(overrides)
}
