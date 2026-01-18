package model

type Skill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	FileName    string `json:"fileName"`
	FilePath    string `json:"filePath"`
	Enabled     bool   `json:"enabled"`
}
