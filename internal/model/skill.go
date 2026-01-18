package model

type Skill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	FileName    string `json:"fileName"`
	FilePath    string `json:"filePath"`
	Enabled     bool   `json:"enabled"`
	Source      string `json:"source"`     // "user" or "plugin"
	PluginName  string `json:"pluginName"` // e.g., "superpowers" (empty for user skills)
}
