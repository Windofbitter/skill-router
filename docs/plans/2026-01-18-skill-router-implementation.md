# Skill Router Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a local web application to manage Claude Code skills with a Go backend and Vue frontend.

**Architecture:** Go backend serves REST API and embedded Vue SPA. Backend scans `~/.claude/commands/` for skills, parses frontmatter, and provides CRUD operations. Frontend displays skills in cards with filtering and search.

**Tech Stack:** Go 1.21+, Vue 3 + Vite + Tailwind CSS, go:embed for SPA bundling

---

## Task 1: Initialize Go Project

**Files:**
- Create: `go.mod`
- Create: `main.go`

**Step 1: Initialize Go module**

```bash
cd /Users/wind/Desktop/projects/skill-router
go mod init github.com/wind/skill-router
```

**Step 2: Create minimal main.go**

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	fmt.Println("Skill Router running at http://localhost:9527")
	http.ListenAndServe(":9527", nil)
}
```

**Step 3: Run and verify**

```bash
go run main.go
# In another terminal:
curl http://localhost:9527/api/health
```
Expected: `{"status":"ok"}`

**Step 4: Commit**

```bash
git add go.mod main.go
git commit -m "feat: initialize go project with health endpoint"
```

---

## Task 2: Implement Frontmatter Parser

**Files:**
- Create: `internal/parser/markdown.go`
- Create: `internal/parser/markdown_test.go`

**Step 1: Write the failing test**

```go
// internal/parser/markdown_test.go
package parser

import (
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	content := `---
name: test-skill
description: A test skill
---

# Content here
`
	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.Name != "test-skill" {
		t.Errorf("expected name 'test-skill', got '%s'", fm.Name)
	}
	if fm.Description != "A test skill" {
		t.Errorf("expected description 'A test skill', got '%s'", fm.Description)
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	content := `# Just markdown without frontmatter`
	fm, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.Name != "" {
		t.Errorf("expected empty name, got '%s'", fm.Name)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/parser/... -v
```
Expected: FAIL - package doesn't exist

**Step 3: Write minimal implementation**

```go
// internal/parser/markdown.go
package parser

import (
	"regexp"
	"strings"
)

type Frontmatter struct {
	Name        string
	Description string
}

var frontmatterRegex = regexp.MustCompile(`(?s)^---\n(.+?)\n---`)

func ParseFrontmatter(content string) (Frontmatter, error) {
	var fm Frontmatter

	matches := frontmatterRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return fm, nil
	}

	lines := strings.Split(matches[1], "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "name":
			fm.Name = value
		case "description":
			fm.Description = value
		}
	}

	return fm, nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/parser/... -v
```
Expected: PASS

**Step 5: Commit**

```bash
git add internal/parser/
git commit -m "feat: add frontmatter parser with tests"
```

---

## Task 3: Implement Skill Service

**Files:**
- Create: `internal/model/skill.go`
- Create: `internal/service/skill.go`
- Create: `internal/service/skill_test.go`

**Step 1: Create skill model**

```go
// internal/model/skill.go
package model

type Skill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	FileName    string `json:"fileName"`
	FilePath    string `json:"filePath"`
	Enabled     bool   `json:"enabled"`
}
```

**Step 2: Write the failing test for ListSkills**

```go
// internal/service/skill_test.go
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
```

**Step 3: Run test to verify it fails**

```bash
go test ./internal/service/... -v
```
Expected: FAIL - package doesn't exist

**Step 4: Write minimal implementation**

```go
// internal/service/skill.go
package service

import (
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
```

**Step 5: Run test to verify it passes**

```bash
go test ./internal/service/... -v
```
Expected: PASS

**Step 6: Commit**

```bash
git add internal/model/ internal/service/
git commit -m "feat: add skill service with list functionality"
```

---

## Task 4: Add Enable/Disable/Delete Methods

**Files:**
- Modify: `internal/service/skill.go`
- Modify: `internal/service/skill_test.go`

**Step 1: Add tests for Enable, Disable, Delete**

```go
// Add to internal/service/skill_test.go

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
```

**Step 2: Run tests to verify they fail**

```bash
go test ./internal/service/... -v
```
Expected: FAIL - methods don't exist

**Step 3: Implement Enable, Disable, Delete**

```go
// Add to internal/service/skill.go

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
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/service/... -v
```
Expected: PASS

**Step 5: Commit**

```bash
git add internal/service/
git commit -m "feat: add enable, disable, delete skill methods"
```

---

## Task 5: Implement HTTP Handlers

**Files:**
- Create: `internal/handler/skill.go`
- Modify: `main.go`

**Step 1: Create handlers**

```go
// internal/handler/skill.go
package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wind/skill-router/internal/service"
)

type SkillHandler struct {
	svc *service.SkillService
}

func NewSkillHandler(svc *service.SkillService) *SkillHandler {
	return &SkillHandler{svc: svc}
}

func (h *SkillHandler) List(w http.ResponseWriter, r *http.Request) {
	skills, err := h.svc.ListSkills()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skills)
}

func (h *SkillHandler) Disable(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	fileName = strings.TrimSuffix(fileName, "/disable")

	if err := h.svc.DisableSkill(fileName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SkillHandler) Enable(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	fileName = strings.TrimSuffix(fileName, "/enable")

	if err := h.svc.EnableSkill(fileName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SkillHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	enabled := r.URL.Query().Get("enabled") == "true"

	if err := h.svc.DeleteSkill(fileName, enabled); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
```

**Step 2: Update main.go**

```go
// main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wind/skill-router/internal/handler"
	"github.com/wind/skill-router/internal/service"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	claudeDir := filepath.Join(homeDir, ".claude")

	svc := service.NewSkillService(claudeDir)
	h := handler.NewSkillHandler(svc)

	http.HandleFunc("/api/skills", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h.List(w, r)
		}
	})

	http.HandleFunc("/api/skills/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/disable") && r.Method == "POST":
			h.Disable(w, r)
		case strings.HasSuffix(path, "/enable") && r.Method == "POST":
			h.Enable(w, r)
		case r.Method == "DELETE":
			h.Delete(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	addr := ":9527"
	url := "http://localhost" + addr

	fmt.Printf("Skill Router running at %s\n", url)

	// Open browser
	go openBrowser(url)

	http.ListenAndServe(addr, nil)
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}
	if cmd != nil {
		cmd.Run()
	}
}
```

**Step 3: Test manually**

```bash
go run .
# In another terminal:
curl http://localhost:9527/api/skills
```
Expected: JSON array of skills

**Step 4: Commit**

```bash
git add internal/handler/ main.go
git commit -m "feat: add HTTP handlers and wire up routes"
```

---

## Task 6: Initialize Vue Frontend

**Files:**
- Create: `web/` directory with Vue project

**Step 1: Create Vue project**

```bash
cd /Users/wind/Desktop/projects/skill-router
npm create vite@latest web -- --template vue-ts
cd web
npm install
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

**Step 2: Configure Tailwind**

```js
// web/tailwind.config.js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

**Step 3: Add Tailwind to CSS**

```css
/* web/src/style.css */
@tailwind base;
@tailwind components;
@tailwind utilities;
```

**Step 4: Configure Vite proxy for development**

```ts
// web/vite.config.ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/api': 'http://localhost:9527'
    }
  }
})
```

**Step 5: Verify it works**

```bash
cd /Users/wind/Desktop/projects/skill-router/web
npm run dev
```
Expected: Vue app running on localhost:5173

**Step 6: Commit**

```bash
cd /Users/wind/Desktop/projects/skill-router
git add web/
git commit -m "feat: initialize vue frontend with tailwind"
```

---

## Task 7: Create API Client

**Files:**
- Create: `web/src/api/skills.ts`
- Create: `web/src/types/skill.ts`

**Step 1: Create types**

```ts
// web/src/types/skill.ts
export interface Skill {
  name: string
  description: string
  fileName: string
  filePath: string
  enabled: boolean
}
```

**Step 2: Create API client**

```ts
// web/src/api/skills.ts
import type { Skill } from '../types/skill'

const API_BASE = '/api'

export async function listSkills(): Promise<Skill[]> {
  const res = await fetch(`${API_BASE}/skills`)
  if (!res.ok) throw new Error('Failed to fetch skills')
  return res.json()
}

export async function disableSkill(fileName: string): Promise<void> {
  const res = await fetch(`${API_BASE}/skills/${fileName}/disable`, {
    method: 'POST'
  })
  if (!res.ok) throw new Error('Failed to disable skill')
}

export async function enableSkill(fileName: string): Promise<void> {
  const res = await fetch(`${API_BASE}/skills/${fileName}/enable`, {
    method: 'POST'
  })
  if (!res.ok) throw new Error('Failed to enable skill')
}

export async function deleteSkill(fileName: string, enabled: boolean): Promise<void> {
  const res = await fetch(`${API_BASE}/skills/${fileName}?enabled=${enabled}`, {
    method: 'DELETE'
  })
  if (!res.ok) throw new Error('Failed to delete skill')
}
```

**Step 3: Commit**

```bash
git add web/src/api/ web/src/types/
git commit -m "feat: add skill api client and types"
```

---

## Task 8: Create SkillCard Component

**Files:**
- Create: `web/src/components/SkillCard.vue`

**Step 1: Create component**

```vue
<!-- web/src/components/SkillCard.vue -->
<script setup lang="ts">
import type { Skill } from '../types/skill'

const props = defineProps<{
  skill: Skill
}>()

const emit = defineEmits<{
  enable: [fileName: string]
  disable: [fileName: string]
  delete: [fileName: string, enabled: boolean]
}>()
</script>

<template>
  <div class="bg-white rounded-lg shadow p-4 border border-gray-200">
    <div class="flex items-start justify-between">
      <div class="flex-1 min-w-0">
        <h3 class="text-lg font-medium text-gray-900 truncate">
          {{ skill.name }}
        </h3>
        <p class="mt-1 text-sm text-gray-500 line-clamp-2">
          {{ skill.description || 'No description' }}
        </p>
      </div>
      <span
        :class="[
          'ml-2 px-2 py-1 text-xs font-medium rounded-full',
          skill.enabled
            ? 'bg-green-100 text-green-800'
            : 'bg-gray-100 text-gray-800'
        ]"
      >
        {{ skill.enabled ? 'Enabled' : 'Disabled' }}
      </span>
    </div>
    <div class="mt-4 flex gap-2">
      <button
        v-if="skill.enabled"
        @click="emit('disable', skill.fileName)"
        class="px-3 py-1 text-sm bg-yellow-100 text-yellow-800 rounded hover:bg-yellow-200"
      >
        Disable
      </button>
      <button
        v-else
        @click="emit('enable', skill.fileName)"
        class="px-3 py-1 text-sm bg-green-100 text-green-800 rounded hover:bg-green-200"
      >
        Enable
      </button>
      <button
        @click="emit('delete', skill.fileName, skill.enabled)"
        class="px-3 py-1 text-sm bg-red-100 text-red-800 rounded hover:bg-red-200"
      >
        Delete
      </button>
    </div>
  </div>
</template>
```

**Step 2: Commit**

```bash
git add web/src/components/SkillCard.vue
git commit -m "feat: add SkillCard component"
```

---

## Task 9: Create Main App View

**Files:**
- Modify: `web/src/App.vue`

**Step 1: Implement main view**

```vue
<!-- web/src/App.vue -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Skill } from './types/skill'
import { listSkills, enableSkill, disableSkill, deleteSkill } from './api/skills'
import SkillCard from './components/SkillCard.vue'

const skills = ref<Skill[]>([])
const filter = ref<'all' | 'enabled' | 'disabled'>('all')
const search = ref('')
const loading = ref(true)

const filteredSkills = computed(() => {
  return skills.value.filter(skill => {
    const matchesFilter =
      filter.value === 'all' ||
      (filter.value === 'enabled' && skill.enabled) ||
      (filter.value === 'disabled' && !skill.enabled)

    const matchesSearch =
      !search.value ||
      skill.name.toLowerCase().includes(search.value.toLowerCase()) ||
      skill.description?.toLowerCase().includes(search.value.toLowerCase())

    return matchesFilter && matchesSearch
  })
})

async function loadSkills() {
  loading.value = true
  try {
    skills.value = await listSkills()
  } finally {
    loading.value = false
  }
}

async function handleEnable(fileName: string) {
  await enableSkill(fileName)
  await loadSkills()
}

async function handleDisable(fileName: string) {
  await disableSkill(fileName)
  await loadSkills()
}

async function handleDelete(fileName: string, enabled: boolean) {
  if (!confirm(`Delete ${fileName}?`)) return
  await deleteSkill(fileName, enabled)
  await loadSkills()
}

onMounted(loadSkills)
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
        <h1 class="text-2xl font-bold text-gray-900">Skill Router</h1>
        <button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          + Add
        </button>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 py-6">
      <div class="flex flex-col sm:flex-row gap-4 mb-6">
        <div class="flex gap-2">
          <button
            v-for="f in ['all', 'enabled', 'disabled'] as const"
            :key="f"
            @click="filter = f"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-medium',
              filter === f
                ? 'bg-blue-600 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-100'
            ]"
          >
            {{ f === 'all' ? 'All' : f === 'enabled' ? 'Enabled' : 'Disabled' }}
          </button>
        </div>
        <input
          v-model="search"
          type="text"
          placeholder="Search skills..."
          class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>

      <div v-if="loading" class="text-center py-12 text-gray-500">
        Loading...
      </div>

      <div v-else-if="filteredSkills.length === 0" class="text-center py-12 text-gray-500">
        No skills found
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <SkillCard
          v-for="skill in filteredSkills"
          :key="skill.filePath"
          :skill="skill"
          @enable="handleEnable"
          @disable="handleDisable"
          @delete="handleDelete"
        />
      </div>
    </main>
  </div>
</template>
```

**Step 2: Verify it works**

```bash
# Terminal 1: Run Go backend
cd /Users/wind/Desktop/projects/skill-router
go run .

# Terminal 2: Run Vue dev server
cd /Users/wind/Desktop/projects/skill-router/web
npm run dev
```
Visit http://localhost:5173 and verify skills are displayed.

**Step 3: Commit**

```bash
git add web/src/App.vue
git commit -m "feat: implement main skill list view"
```

---

## Task 10: Add Upload Skill Feature

**Files:**
- Modify: `internal/handler/skill.go`
- Modify: `internal/service/skill.go`
- Create: `web/src/components/AddSkillModal.vue`

**Step 1: Add upload endpoint to Go backend**

```go
// Add to internal/service/skill.go

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
```

```go
// Add to internal/handler/skill.go

func (h *SkillHandler) Upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	overwrite := r.FormValue("overwrite") == "true"

	if err := h.svc.SaveSkill(header.Filename, content, overwrite); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
```

Add import `"io"` and `"fmt"` at the top of respective files.

**Step 2: Add route in main.go**

```go
// Add to main.go in the http handlers section

http.HandleFunc("/api/skills/upload", func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.Upload(w, r)
	}
})
```

**Step 3: Add upload to frontend API**

```ts
// Add to web/src/api/skills.ts

export async function uploadSkill(file: File, overwrite: boolean = false): Promise<void> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('overwrite', String(overwrite))

  const res = await fetch(`${API_BASE}/skills/upload`, {
    method: 'POST',
    body: formData
  })
  if (!res.ok) {
    if (res.status === 409) throw new Error('File already exists')
    throw new Error('Failed to upload skill')
  }
}
```

**Step 4: Create AddSkillModal component**

```vue
<!-- web/src/components/AddSkillModal.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import { uploadSkill } from '../api/skills'

const emit = defineEmits<{
  close: []
  added: []
}>()

const dragOver = ref(false)
const uploading = ref(false)
const error = ref('')

async function handleFiles(files: FileList | null) {
  if (!files || files.length === 0) return

  const file = files[0]
  if (!file.name.endsWith('.md')) {
    error.value = 'Only .md files are allowed'
    return
  }

  uploading.value = true
  error.value = ''

  try {
    await uploadSkill(file)
    emit('added')
    emit('close')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Upload failed'
  } finally {
    uploading.value = false
  }
}

function onDrop(e: DragEvent) {
  dragOver.value = false
  handleFiles(e.dataTransfer?.files ?? null)
}

function onFileSelect(e: Event) {
  const input = e.target as HTMLInputElement
  handleFiles(input.files)
}
</script>

<template>
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
    <div class="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-bold">Add Skill</h2>
        <button @click="emit('close')" class="text-gray-500 hover:text-gray-700">
          &times;
        </button>
      </div>

      <div
        @dragover.prevent="dragOver = true"
        @dragleave="dragOver = false"
        @drop.prevent="onDrop"
        :class="[
          'border-2 border-dashed rounded-lg p-8 text-center transition-colors',
          dragOver ? 'border-blue-500 bg-blue-50' : 'border-gray-300'
        ]"
      >
        <p class="text-gray-600 mb-2">Drag and drop a .md file here</p>
        <p class="text-gray-400 text-sm mb-4">or</p>
        <label class="px-4 py-2 bg-blue-600 text-white rounded-lg cursor-pointer hover:bg-blue-700">
          Choose File
          <input type="file" accept=".md" class="hidden" @change="onFileSelect" />
        </label>
      </div>

      <p v-if="error" class="mt-4 text-red-600 text-sm">{{ error }}</p>
      <p v-if="uploading" class="mt-4 text-gray-600 text-sm">Uploading...</p>
    </div>
  </div>
</template>
```

**Step 5: Integrate modal into App.vue**

```vue
<!-- Add to web/src/App.vue -->
<!-- In script setup, add: -->
const showAddModal = ref(false)

<!-- Update the Add button: -->
<button @click="showAddModal = true" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
  + Add
</button>

<!-- Add modal at the end of template, before closing </div>: -->
<AddSkillModal
  v-if="showAddModal"
  @close="showAddModal = false"
  @added="loadSkills"
/>

<!-- Add import: -->
import AddSkillModal from './components/AddSkillModal.vue'
```

**Step 6: Test upload**

Run both servers and test uploading a .md file.

**Step 7: Commit**

```bash
git add .
git commit -m "feat: add skill upload functionality"
```

---

## Task 11: Add GitHub Install Feature

**Files:**
- Create: `internal/github/client.go`
- Modify: `internal/handler/skill.go`
- Modify: `web/src/components/AddSkillModal.vue`

**Step 1: Create GitHub client**

```go
// internal/github/client.go
package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type File struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
}

var repoURLRegex = regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`)

func FetchSkillFiles(repoURL string) ([]File, error) {
	matches := repoURLRegex.FindStringSubmatch(repoURL)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid github URL")
	}

	owner, repo := matches[1], matches[2]
	repo = strings.TrimSuffix(repo, ".git")

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/.claude/commands", owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("no .claude/commands directory found")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github API error: %d", resp.StatusCode)
	}

	var files []File
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	var mdFiles []File
	for _, f := range files {
		if f.Type == "file" && strings.HasSuffix(f.Name, ".md") {
			mdFiles = append(mdFiles, f)
		}
	}

	return mdFiles, nil
}

func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
```

**Step 2: Add install handler**

```go
// Add to internal/handler/skill.go

type InstallRequest struct {
	URL string `json:"url"`
}

func (h *SkillHandler) Install(w http.ResponseWriter, r *http.Request) {
	var req InstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	files, err := github.FetchSkillFiles(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	installed := 0
	for _, f := range files {
		content, err := github.DownloadFile(f.DownloadURL)
		if err != nil {
			continue
		}
		if err := h.svc.SaveSkill(f.Name, content, false); err != nil {
			continue
		}
		installed++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"installed": installed})
}
```

Add import `"github.com/wind/skill-router/internal/github"` to handler.

**Step 3: Add route in main.go**

```go
http.HandleFunc("/api/skills/install", func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.Install(w, r)
	}
})
```

**Step 4: Add install to frontend API**

```ts
// Add to web/src/api/skills.ts

export async function installFromGithub(url: string): Promise<{ installed: number }> {
  const res = await fetch(`${API_BASE}/skills/install`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url })
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Failed to install skills')
  }
  return res.json()
}
```

**Step 5: Update AddSkillModal with GitHub tab**

```vue
<!-- web/src/components/AddSkillModal.vue - Replace entire content -->
<script setup lang="ts">
import { ref } from 'vue'
import { uploadSkill, installFromGithub } from '../api/skills'

const emit = defineEmits<{
  close: []
  added: []
}>()

const tab = ref<'upload' | 'github'>('upload')
const dragOver = ref(false)
const loading = ref(false)
const error = ref('')
const githubUrl = ref('')

async function handleFiles(files: FileList | null) {
  if (!files || files.length === 0) return

  const file = files[0]
  if (!file.name.endsWith('.md')) {
    error.value = 'Only .md files are allowed'
    return
  }

  loading.value = true
  error.value = ''

  try {
    await uploadSkill(file)
    emit('added')
    emit('close')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Upload failed'
  } finally {
    loading.value = false
  }
}

async function handleGithubInstall() {
  if (!githubUrl.value) return

  loading.value = true
  error.value = ''

  try {
    const result = await installFromGithub(githubUrl.value)
    if (result.installed === 0) {
      error.value = 'No skills found in repository'
    } else {
      emit('added')
      emit('close')
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Install failed'
  } finally {
    loading.value = false
  }
}

function onDrop(e: DragEvent) {
  dragOver.value = false
  handleFiles(e.dataTransfer?.files ?? null)
}

function onFileSelect(e: Event) {
  const input = e.target as HTMLInputElement
  handleFiles(input.files)
}
</script>

<template>
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
    <div class="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-bold">Add Skill</h2>
        <button @click="emit('close')" class="text-gray-500 hover:text-gray-700 text-2xl">
          &times;
        </button>
      </div>

      <div class="flex gap-2 mb-4">
        <button
          @click="tab = 'upload'"
          :class="[
            'px-4 py-2 rounded-lg text-sm font-medium flex-1',
            tab === 'upload' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700'
          ]"
        >
          Upload File
        </button>
        <button
          @click="tab = 'github'"
          :class="[
            'px-4 py-2 rounded-lg text-sm font-medium flex-1',
            tab === 'github' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700'
          ]"
        >
          From GitHub
        </button>
      </div>

      <div v-if="tab === 'upload'">
        <div
          @dragover.prevent="dragOver = true"
          @dragleave="dragOver = false"
          @drop.prevent="onDrop"
          :class="[
            'border-2 border-dashed rounded-lg p-8 text-center transition-colors',
            dragOver ? 'border-blue-500 bg-blue-50' : 'border-gray-300'
          ]"
        >
          <p class="text-gray-600 mb-2">Drag and drop a .md file here</p>
          <p class="text-gray-400 text-sm mb-4">or</p>
          <label class="px-4 py-2 bg-blue-600 text-white rounded-lg cursor-pointer hover:bg-blue-700">
            Choose File
            <input type="file" accept=".md" class="hidden" @change="onFileSelect" />
          </label>
        </div>
      </div>

      <div v-else>
        <input
          v-model="githubUrl"
          type="text"
          placeholder="https://github.com/user/repo"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <p class="mt-2 text-sm text-gray-500">
          Will install all skills from .claude/commands/
        </p>
        <button
          @click="handleGithubInstall"
          :disabled="!githubUrl || loading"
          class="mt-4 w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
        >
          Install
        </button>
      </div>

      <p v-if="error" class="mt-4 text-red-600 text-sm">{{ error }}</p>
      <p v-if="loading" class="mt-4 text-gray-600 text-sm">Loading...</p>
    </div>
  </div>
</template>
```

**Step 6: Test GitHub install**

Test with a real GitHub repo that has `.claude/commands/` directory.

**Step 7: Commit**

```bash
git add .
git commit -m "feat: add github skill installation"
```

---

## Task 12: Embed Frontend and Build Binary

**Files:**
- Create: `embed.go`
- Create: `Makefile`
- Modify: `main.go`

**Step 1: Create embed.go**

```go
// embed.go
package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/dist/*
var webFS embed.FS

func getFileSystem() http.FileSystem {
	subFS, _ := fs.Sub(webFS, "web/dist")
	return http.FS(subFS)
}
```

**Step 2: Update main.go to serve embedded files**

```go
// Add to main.go, after API routes

// Serve static files
fileServer := http.FileServer(getFileSystem())
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// Try to serve the file
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Check if file exists in embedded FS
	if f, err := webFS.Open("web/dist" + path); err == nil {
		f.Close()
		fileServer.ServeHTTP(w, r)
		return
	}

	// Fallback to index.html for SPA routing
	r.URL.Path = "/"
	fileServer.ServeHTTP(w, r)
})
```

**Step 3: Create Makefile**

```makefile
# Makefile
.PHONY: build dev clean

build: build-frontend build-backend

build-frontend:
	cd web && npm install && npm run build

build-backend:
	go build -o skill-router .

dev:
	@echo "Run 'go run .' in one terminal"
	@echo "Run 'cd web && npm run dev' in another terminal"

clean:
	rm -rf skill-router web/dist web/node_modules
```

**Step 4: Build and test**

```bash
make build
./skill-router
```
Expected: Opens browser with working app

**Step 5: Commit**

```bash
git add embed.go Makefile main.go
git commit -m "feat: embed frontend and create build system"
```

---

## Task 13: Final Polish and README

**Files:**
- Create: `README.md`

**Step 1: Create README**

```markdown
# Skill Router

A local web application to manage Claude Code skills.

## Features

- View all skills from `~/.claude/commands/`
- Enable/disable skills (moves to `~/.claude/skills-disabled/`)
- Delete skills
- Upload .md skill files
- Install skills from GitHub repositories

## Installation

Download the binary for your platform from [Releases](https://github.com/wind/skill-router/releases).

Or build from source:

```bash
make build
```

## Usage

```bash
./skill-router
```

This starts the server and opens your browser to http://localhost:9527

## Development

```bash
# Terminal 1: Run Go backend
go run .

# Terminal 2: Run Vue dev server
cd web && npm run dev
```

## License

MIT
```

**Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README"
```

---

## Summary

| Task | Description |
|------|-------------|
| 1 | Initialize Go project |
| 2 | Implement frontmatter parser |
| 3 | Implement skill service (list) |
| 4 | Add enable/disable/delete methods |
| 5 | Implement HTTP handlers |
| 6 | Initialize Vue frontend |
| 7 | Create API client |
| 8 | Create SkillCard component |
| 9 | Create main app view |
| 10 | Add upload feature |
| 11 | Add GitHub install feature |
| 12 | Embed frontend and build binary |
| 13 | Final polish and README |
