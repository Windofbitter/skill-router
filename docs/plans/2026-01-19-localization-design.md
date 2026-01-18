# Localization and README Design

## Overview

Add Chinese localization to the Skill Router UI and create comprehensive README documentation in both English and Chinese.

## Localization Architecture

### Library & Setup

Using `vue-i18n@10` (latest for Vue 3) with the Composition API.

```
web/src/
├── i18n/
│   ├── index.ts          # vue-i18n setup & language detection
│   ├── locales/
│   │   ├── en.json       # English translations
│   │   └── zh.json       # Chinese translations
```

### Language Detection Logic

1. Check `localStorage` for user's saved preference
2. If none, detect from `navigator.language` (e.g., `zh-CN` → Chinese)
3. Fall back to English if language not supported

### Language Toggle

A small language switcher in the header (top-right corner), showing current language code (EN/中) that cycles between options on click. Selection saves to `localStorage`.

### String Organization

Translations structured by component/feature:

```json
{
  "header": { "title": "Skill Router", "add": "Add" },
  "filters": { "all": "All", "user": "User", ... },
  "skillCard": { "enable": "Enable", "disable": "Disable", ... },
  "errors": { "fetchFailed": "Failed to fetch skills", ... }
}
```

## Component Changes

### Files to Modify

| File | Changes |
|------|---------|
| `web/package.json` | Add `vue-i18n` dependency |
| `web/src/main.ts` | Import and register i18n plugin |
| `web/src/App.vue` | Replace hardcoded text with `$t()` calls, add language toggle |
| `web/src/components/AddSkillModal.vue` | Translate labels, placeholders, error messages |
| `web/src/components/PluginGroup.vue` | Translate buttons and skill count text |
| `web/src/components/SkillCard.vue` | Translate status badges and action buttons |
| `web/src/api/skills.ts` | Export error keys instead of hardcoded strings, let components translate |

### Translation Count

- ~15 UI labels (buttons, headers, tabs)
- ~8 status/state messages (Loading, No skills found, etc.)
- ~12 error messages (from API layer)
- ~5 confirmation dialogs
- **Total: ~40 translation keys**

### Dynamic Content

For messages with variables like "Delete {fileName}?", use vue-i18n interpolation:

```js
$t('confirm.deleteSkill', { fileName: skill.name })
```

## README Documentation

### File Structure

```
/README.md           # English (primary)
/README_CN.md        # Chinese (full translation)
/docs/images/        # Screenshots shared by both
```

### Content Outline

1. Header - Project name, brief tagline
2. Screenshot - Main UI showing skill list
3. Features - Bullet list of capabilities
4. Installation - Download binary or build from source
5. Usage - How to run, access the web UI
6. Development - Build instructions, project structure
7. License - MIT

### Screenshots to Capture

| Screenshot | Description |
|------------|-------------|
| `main-view.png` | Main interface with skills listed |
| `add-skill-modal.png` | The Add Skill modal (upload/GitHub tabs) |

## Implementation Order

### Phase 1: i18n Setup
1. Install vue-i18n
2. Create i18n config with language detection
3. Create English translation file (extract all strings)
4. Create Chinese translation file

### Phase 2: Component Updates
5. Wire i18n plugin to Vue app
6. Update App.vue with translations + language toggle
7. Update AddSkillModal.vue
8. Update PluginGroup.vue
9. Update SkillCard.vue
10. Update api/skills.ts error handling

### Phase 3: Documentation
11. Start servers, capture screenshots with agent-browser
12. Write README.md (English)
13. Write README_CN.md (Chinese)
14. Commit all changes
