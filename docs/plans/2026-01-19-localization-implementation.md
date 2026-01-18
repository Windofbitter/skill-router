# Localization & README Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Chinese localization to the Vue UI and create bilingual README documentation with screenshots.

**Architecture:** Vue-i18n for reactive translations with browser auto-detect + manual toggle. JSON translation files organized by component. READMEs as full mirrors with shared screenshots.

**Tech Stack:** Vue 3, vue-i18n@10, TypeScript, Vite, Tailwind CSS

---

## Phase 1: i18n Infrastructure

### Task 1: Install vue-i18n

**Files:**
- Modify: `web/package.json`

**Step 1: Add vue-i18n dependency**

Run:
```bash
cd /Users/wind/Desktop/projects/skill-router/web && npm install vue-i18n@10
```

**Step 2: Verify installation**

Run:
```bash
cd /Users/wind/Desktop/projects/skill-router/web && npm ls vue-i18n
```
Expected: `vue-i18n@10.x.x`

---

### Task 2: Create English translation file

**Files:**
- Create: `web/src/i18n/locales/en.json`

**Step 1: Create the English translations file**

Create `web/src/i18n/locales/en.json`:
```json
{
  "header": {
    "title": "Skill Router",
    "add": "+ Add"
  },
  "filters": {
    "all": "All",
    "user": "User",
    "plugins": "Plugins",
    "enabled": "Enabled",
    "disabled": "Disabled",
    "searchPlaceholder": "Search skills..."
  },
  "status": {
    "loading": "Loading...",
    "noSkills": "No skills found",
    "enabled": "Enabled",
    "disabled": "Disabled"
  },
  "sections": {
    "userSkills": "User Skills",
    "pluginSkills": "Plugin Skills"
  },
  "skillCard": {
    "noDescription": "No description",
    "enable": "Enable",
    "disable": "Disable",
    "delete": "Delete"
  },
  "pluginGroup": {
    "skillsEnabled": "{count}/{total} skills enabled",
    "enableAll": "Enable All",
    "disableAll": "Disable All",
    "remove": "Remove"
  },
  "addModal": {
    "title": "Add Skill",
    "uploadFile": "Upload File",
    "fromGithub": "From GitHub",
    "dragDrop": "Drag and drop a .md file here",
    "or": "or",
    "chooseFile": "Choose File",
    "githubPlaceholder": "https://github.com/user/repo",
    "githubHelper": "Will install all skills from .claude/commands/",
    "install": "Install"
  },
  "confirm": {
    "deleteSkill": "Delete {fileName}?",
    "removePlugin": "Remove plugin \"{pluginName}\"? This will delete all its skills from the cache."
  },
  "errors": {
    "fetchFailed": "Failed to fetch skills",
    "disableFailed": "Failed to disable skill",
    "enableFailed": "Failed to enable skill",
    "deleteFailed": "Failed to delete skill",
    "fileExists": "File already exists",
    "uploadFailed": "Failed to upload skill",
    "installFailed": "Failed to install skills",
    "disablePluginSkillFailed": "Failed to disable plugin skill",
    "enablePluginSkillFailed": "Failed to enable plugin skill",
    "disablePluginFailed": "Failed to disable plugin",
    "enablePluginFailed": "Failed to enable plugin",
    "deletePluginFailed": "Failed to delete plugin",
    "onlyMdAllowed": "Only .md files are allowed",
    "noSkillsInRepo": "No skills found in repository"
  },
  "language": {
    "en": "EN",
    "zh": "中"
  }
}
```

---

### Task 3: Create Chinese translation file

**Files:**
- Create: `web/src/i18n/locales/zh.json`

**Step 1: Create the Chinese translations file**

Create `web/src/i18n/locales/zh.json`:
```json
{
  "header": {
    "title": "技能路由",
    "add": "+ 添加"
  },
  "filters": {
    "all": "全部",
    "user": "用户",
    "plugins": "插件",
    "enabled": "已启用",
    "disabled": "已禁用",
    "searchPlaceholder": "搜索技能..."
  },
  "status": {
    "loading": "加载中...",
    "noSkills": "未找到技能",
    "enabled": "已启用",
    "disabled": "已禁用"
  },
  "sections": {
    "userSkills": "用户技能",
    "pluginSkills": "插件技能"
  },
  "skillCard": {
    "noDescription": "暂无描述",
    "enable": "启用",
    "disable": "禁用",
    "delete": "删除"
  },
  "pluginGroup": {
    "skillsEnabled": "{count}/{total} 个技能已启用",
    "enableAll": "全部启用",
    "disableAll": "全部禁用",
    "remove": "移除"
  },
  "addModal": {
    "title": "添加技能",
    "uploadFile": "上传文件",
    "fromGithub": "从 GitHub",
    "dragDrop": "拖放 .md 文件到此处",
    "or": "或",
    "chooseFile": "选择文件",
    "githubPlaceholder": "https://github.com/user/repo",
    "githubHelper": "将从 .claude/commands/ 安装所有技能",
    "install": "安装"
  },
  "confirm": {
    "deleteSkill": "删除 {fileName}？",
    "removePlugin": "移除插件 \"{pluginName}\"？这将从缓存中删除其所有技能。"
  },
  "errors": {
    "fetchFailed": "获取技能失败",
    "disableFailed": "禁用技能失败",
    "enableFailed": "启用技能失败",
    "deleteFailed": "删除技能失败",
    "fileExists": "文件已存在",
    "uploadFailed": "上传技能失败",
    "installFailed": "安装技能失败",
    "disablePluginSkillFailed": "禁用插件技能失败",
    "enablePluginSkillFailed": "启用插件技能失败",
    "disablePluginFailed": "禁用插件失败",
    "enablePluginFailed": "启用插件失败",
    "deletePluginFailed": "删除插件失败",
    "onlyMdAllowed": "仅允许 .md 文件",
    "noSkillsInRepo": "仓库中未找到技能"
  },
  "language": {
    "en": "EN",
    "zh": "中"
  }
}
```

---

### Task 4: Create i18n configuration

**Files:**
- Create: `web/src/i18n/index.ts`

**Step 1: Create the i18n setup file**

Create `web/src/i18n/index.ts`:
```typescript
import { createI18n } from 'vue-i18n'
import en from './locales/en.json'
import zh from './locales/zh.json'

type MessageSchema = typeof en

const STORAGE_KEY = 'skill-router-locale'

function detectLocale(): 'en' | 'zh' {
  // Check localStorage first
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved === 'en' || saved === 'zh') {
    return saved
  }

  // Detect from browser
  const browserLang = navigator.language.toLowerCase()
  if (browserLang.startsWith('zh')) {
    return 'zh'
  }

  return 'en'
}

export function saveLocale(locale: 'en' | 'zh') {
  localStorage.setItem(STORAGE_KEY, locale)
}

export const i18n = createI18n<[MessageSchema], 'en' | 'zh'>({
  legacy: false,
  locale: detectLocale(),
  fallbackLocale: 'en',
  messages: {
    en,
    zh
  }
})
```

---

### Task 5: Wire i18n to Vue app

**Files:**
- Modify: `web/src/main.ts`

**Step 1: Update main.ts to use i18n**

Replace `web/src/main.ts` content:
```typescript
import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { i18n } from './i18n'

createApp(App).use(i18n).mount('#app')
```

**Step 2: Commit i18n infrastructure**

```bash
git add web/package.json web/package-lock.json web/src/i18n/ web/src/main.ts
git commit -m "feat: add vue-i18n infrastructure with en/zh translations"
```

---

## Phase 2: Component Updates

### Task 6: Update App.vue with translations

**Files:**
- Modify: `web/src/App.vue`

**Step 1: Update App.vue**

Replace `web/src/App.vue` content:
```vue
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Skill } from './types/skill'
import {
  listSkills,
  enableSkill,
  disableSkill,
  deleteSkill,
  enablePluginSkill,
  disablePluginSkill,
  enablePlugin,
  disablePlugin,
  deletePlugin
} from './api/skills'
import SkillCard from './components/SkillCard.vue'
import PluginGroup from './components/PluginGroup.vue'
import AddSkillModal from './components/AddSkillModal.vue'
import { saveLocale } from './i18n'

const { t, locale } = useI18n()

const skills = ref<Skill[]>([])
const filter = ref<'all' | 'enabled' | 'disabled'>('all')
const sourceFilter = ref<'all' | 'user' | 'plugin'>('all')
const search = ref('')
const loading = ref(true)
const showAddModal = ref(false)

function toggleLocale() {
  const newLocale = locale.value === 'en' ? 'zh' : 'en'
  locale.value = newLocale
  saveLocale(newLocale)
}

const filteredSkills = computed(() => {
  return skills.value.filter(skill => {
    const matchesFilter =
      filter.value === 'all' ||
      (filter.value === 'enabled' && skill.enabled) ||
      (filter.value === 'disabled' && !skill.enabled)

    const matchesSource =
      sourceFilter.value === 'all' ||
      skill.source === sourceFilter.value

    const matchesSearch =
      !search.value ||
      skill.name.toLowerCase().includes(search.value.toLowerCase()) ||
      skill.description?.toLowerCase().includes(search.value.toLowerCase()) ||
      skill.pluginName?.toLowerCase().includes(search.value.toLowerCase())

    return matchesFilter && matchesSource && matchesSearch
  })
})

const userSkills = computed(() => filteredSkills.value.filter(s => s.source === 'user'))
const pluginSkills = computed(() => filteredSkills.value.filter(s => s.source === 'plugin'))

// Group plugin skills by plugin name
const pluginGroups = computed(() => {
  const groups: Record<string, Skill[]> = {}
  for (const skill of pluginSkills.value) {
    const key = skill.pluginName || 'unknown'
    if (!groups[key]) {
      groups[key] = []
    }
    groups[key].push(skill)
  }
  return groups
})

const userCount = computed(() => skills.value.filter(s => s.source === 'user').length)
const pluginGroupCount = computed(() => {
  const names = new Set(skills.value.filter(s => s.source === 'plugin').map(s => s.pluginName))
  return names.size
})

async function loadSkills() {
  loading.value = true
  try {
    skills.value = await listSkills() || []
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
  if (!confirm(t('confirm.deleteSkill', { fileName }))) return
  await deleteSkill(fileName, enabled)
  await loadSkills()
}

async function handleEnablePluginSkill(pluginName: string, skillName: string) {
  await enablePluginSkill(pluginName, skillName)
  await loadSkills()
}

async function handleDisablePluginSkill(pluginName: string, skillName: string) {
  await disablePluginSkill(pluginName, skillName)
  await loadSkills()
}

async function handleEnablePlugin(pluginName: string) {
  await enablePlugin(pluginName)
  await loadSkills()
}

async function handleDisablePlugin(pluginName: string) {
  await disablePlugin(pluginName)
  await loadSkills()
}

async function handleDeletePlugin(pluginName: string) {
  if (!confirm(t('confirm.removePlugin', { pluginName }))) return
  await deletePlugin(pluginName)
  await loadSkills()
}

onMounted(loadSkills)
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
        <h1 class="text-2xl font-bold text-gray-900">{{ t('header.title') }}</h1>
        <div class="flex items-center gap-3">
          <button
            @click="toggleLocale"
            class="px-3 py-1.5 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200"
          >
            {{ locale === 'en' ? t('language.zh') : t('language.en') }}
          </button>
          <button @click="showAddModal = true" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
            {{ t('header.add') }}
          </button>
        </div>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 py-6">
      <!-- Filters -->
      <div class="flex flex-col gap-4 mb-6">
        <!-- Source filter -->
        <div class="flex gap-2">
          <button
            @click="sourceFilter = 'all'"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-medium',
              sourceFilter === 'all'
                ? 'bg-gray-800 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-100'
            ]"
          >
            {{ t('filters.all') }} ({{ skills.length }})
          </button>
          <button
            @click="sourceFilter = 'user'"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-medium',
              sourceFilter === 'user'
                ? 'bg-gray-800 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-100'
            ]"
          >
            {{ t('filters.user') }} ({{ userCount }})
          </button>
          <button
            @click="sourceFilter = 'plugin'"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-medium',
              sourceFilter === 'plugin'
                ? 'bg-purple-600 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-100'
            ]"
          >
            {{ t('filters.plugins') }} ({{ pluginGroupCount }})
          </button>
        </div>

        <!-- Status filter + search -->
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="flex gap-2">
            <button
              @click="filter = 'all'"
              :class="[
                'px-4 py-2 rounded-lg text-sm font-medium',
                filter === 'all'
                  ? 'bg-blue-600 text-white'
                  : 'bg-white text-gray-700 hover:bg-gray-100'
              ]"
            >
              {{ t('filters.all') }}
            </button>
            <button
              @click="filter = 'enabled'"
              :class="[
                'px-4 py-2 rounded-lg text-sm font-medium',
                filter === 'enabled'
                  ? 'bg-blue-600 text-white'
                  : 'bg-white text-gray-700 hover:bg-gray-100'
              ]"
            >
              {{ t('filters.enabled') }}
            </button>
            <button
              @click="filter = 'disabled'"
              :class="[
                'px-4 py-2 rounded-lg text-sm font-medium',
                filter === 'disabled'
                  ? 'bg-blue-600 text-white'
                  : 'bg-white text-gray-700 hover:bg-gray-100'
              ]"
            >
              {{ t('filters.disabled') }}
            </button>
          </div>
          <input
            v-model="search"
            type="text"
            :placeholder="t('filters.searchPlaceholder')"
            class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
      </div>

      <div v-if="loading" class="text-center py-12 text-gray-500">
        {{ t('status.loading') }}
      </div>

      <div v-else-if="filteredSkills.length === 0" class="text-center py-12 text-gray-500">
        {{ t('status.noSkills') }}
      </div>

      <!-- Plugin view: grouped by plugin -->
      <div v-else-if="sourceFilter === 'plugin'" class="space-y-4">
        <PluginGroup
          v-for="(groupSkills, pluginName) in pluginGroups"
          :key="pluginName"
          :plugin-name="pluginName"
          :skills="groupSkills"
          @enable-plugin="handleEnablePlugin"
          @disable-plugin="handleDisablePlugin"
          @delete-plugin="handleDeletePlugin"
          @enable-skill="handleEnablePluginSkill"
          @disable-skill="handleDisablePluginSkill"
        />
      </div>

      <!-- User/All view: card grid -->
      <div v-else class="space-y-6">
        <!-- User skills -->
        <div v-if="userSkills.length > 0">
          <h2 v-if="sourceFilter === 'all'" class="text-lg font-semibold text-gray-700 mb-3">{{ t('sections.userSkills') }}</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <SkillCard
              v-for="skill in userSkills"
              :key="skill.filePath"
              :skill="skill"
              @enable="handleEnable"
              @disable="handleDisable"
              @delete="handleDelete"
            />
          </div>
        </div>

        <!-- Plugin skills (when viewing "all") -->
        <div v-if="sourceFilter === 'all' && Object.keys(pluginGroups).length > 0">
          <h2 class="text-lg font-semibold text-gray-700 mb-3">{{ t('sections.pluginSkills') }}</h2>
          <div class="space-y-4">
            <PluginGroup
              v-for="(groupSkills, pluginName) in pluginGroups"
              :key="pluginName"
              :plugin-name="pluginName"
              :skills="groupSkills"
              @enable-plugin="handleEnablePlugin"
              @disable-plugin="handleDisablePlugin"
              @delete-plugin="handleDeletePlugin"
              @enable-skill="handleEnablePluginSkill"
              @disable-skill="handleDisablePluginSkill"
            />
          </div>
        </div>
      </div>
    </main>

    <AddSkillModal
      v-if="showAddModal"
      @close="showAddModal = false"
      @added="loadSkills"
    />
  </div>
</template>
```

---

### Task 7: Update SkillCard.vue with translations

**Files:**
- Modify: `web/src/components/SkillCard.vue`

**Step 1: Update SkillCard.vue**

Replace `web/src/components/SkillCard.vue` content:
```vue
<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Skill } from '../types/skill'

const { t } = useI18n()

const props = defineProps<{
  skill: Skill
}>()

const emit = defineEmits<{
  enable: [fileName: string]
  disable: [fileName: string]
  delete: [fileName: string, enabled: boolean]
  enablePlugin: [pluginName: string, skillName: string]
  disablePlugin: [pluginName: string, skillName: string]
}>()
</script>

<template>
  <div class="bg-white rounded-lg shadow p-4 border border-gray-200">
    <div class="flex items-start justify-between">
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2">
          <h3 class="text-lg font-medium text-gray-900 truncate">
            {{ skill.name }}
          </h3>
          <span
            v-if="skill.source === 'plugin'"
            class="px-2 py-0.5 text-xs font-medium rounded bg-purple-100 text-purple-800"
          >
            {{ skill.pluginName }}
          </span>
        </div>
        <p class="mt-1 text-sm text-gray-500 line-clamp-2">
          {{ skill.description || t('skillCard.noDescription') }}
        </p>
      </div>
      <span
        :class="[
          'ml-2 px-2 py-1 text-xs font-medium rounded-full shrink-0',
          skill.enabled
            ? 'bg-green-100 text-green-800'
            : 'bg-gray-100 text-gray-800'
        ]"
      >
        {{ skill.enabled ? t('status.enabled') : t('status.disabled') }}
      </span>
    </div>

    <!-- User skills: full control -->
    <div v-if="skill.source === 'user'" class="mt-4 flex gap-2">
      <button
        v-if="skill.enabled"
        @click="emit('disable', skill.fileName)"
        class="px-3 py-1 text-sm bg-yellow-100 text-yellow-800 rounded hover:bg-yellow-200"
      >
        {{ t('skillCard.disable') }}
      </button>
      <button
        v-else
        @click="emit('enable', skill.fileName)"
        class="px-3 py-1 text-sm bg-green-100 text-green-800 rounded hover:bg-green-200"
      >
        {{ t('skillCard.enable') }}
      </button>
      <button
        @click="emit('delete', skill.fileName, skill.enabled)"
        class="px-3 py-1 text-sm bg-red-100 text-red-800 rounded hover:bg-red-200"
      >
        {{ t('skillCard.delete') }}
      </button>
    </div>

    <!-- Plugin skills: toggle only -->
    <div v-else class="mt-4 flex items-center gap-2">
      <button
        v-if="skill.enabled"
        @click="emit('disablePlugin', skill.pluginName, skill.fileName)"
        class="px-3 py-1 text-sm bg-yellow-100 text-yellow-800 rounded hover:bg-yellow-200"
      >
        {{ t('skillCard.disable') }}
      </button>
      <button
        v-else
        @click="emit('enablePlugin', skill.pluginName, skill.fileName)"
        class="px-3 py-1 text-sm bg-green-100 text-green-800 rounded hover:bg-green-200"
      >
        {{ t('skillCard.enable') }}
      </button>
    </div>
  </div>
</template>
```

---

### Task 8: Update PluginGroup.vue with translations

**Files:**
- Modify: `web/src/components/PluginGroup.vue`

**Step 1: Update PluginGroup.vue**

Replace `web/src/components/PluginGroup.vue` content:
```vue
<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Skill } from '../types/skill'
import SkillCard from './SkillCard.vue'

const { t } = useI18n()

const props = defineProps<{
  pluginName: string
  skills: Skill[]
}>()

const emit = defineEmits<{
  enablePlugin: [pluginName: string]
  disablePlugin: [pluginName: string]
  deletePlugin: [pluginName: string]
  enableSkill: [pluginName: string, skillName: string]
  disableSkill: [pluginName: string, skillName: string]
}>()

const expanded = ref(true)

const allEnabled = computed(() => props.skills.every(s => s.enabled))
const allDisabled = computed(() => props.skills.every(s => !s.enabled))
const enabledCount = computed(() => props.skills.filter(s => s.enabled).length)
</script>

<template>
  <div class="bg-white rounded-lg shadow border border-purple-200 overflow-hidden">
    <!-- Plugin Header -->
    <div class="bg-purple-50 px-4 py-3 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button
          @click="expanded = !expanded"
          class="text-purple-600 hover:text-purple-800"
        >
          <svg
            class="w-5 h-5 transition-transform"
            :class="{ 'rotate-90': expanded }"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
        <div>
          <h3 class="font-semibold text-purple-900">{{ pluginName }}</h3>
          <p class="text-xs text-purple-600">
            {{ t('pluginGroup.skillsEnabled', { count: enabledCount, total: skills.length }) }}
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          v-if="!allEnabled"
          @click="emit('enablePlugin', pluginName)"
          class="px-3 py-1 text-sm bg-green-100 text-green-800 rounded hover:bg-green-200"
        >
          {{ t('pluginGroup.enableAll') }}
        </button>
        <button
          v-if="!allDisabled"
          @click="emit('disablePlugin', pluginName)"
          class="px-3 py-1 text-sm bg-yellow-100 text-yellow-800 rounded hover:bg-yellow-200"
        >
          {{ t('pluginGroup.disableAll') }}
        </button>
        <button
          @click="emit('deletePlugin', pluginName)"
          class="px-3 py-1 text-sm bg-red-100 text-red-800 rounded hover:bg-red-200"
        >
          {{ t('pluginGroup.remove') }}
        </button>
      </div>
    </div>

    <!-- Skills List -->
    <div v-if="expanded" class="p-4 grid grid-cols-1 sm:grid-cols-2 gap-3">
      <SkillCard
        v-for="skill in skills"
        :key="skill.filePath"
        :skill="skill"
        @enablePlugin="emit('enableSkill', $event, skill.fileName)"
        @disablePlugin="emit('disableSkill', $event, skill.fileName)"
      />
    </div>
  </div>
</template>
```

---

### Task 9: Update AddSkillModal.vue with translations

**Files:**
- Modify: `web/src/components/AddSkillModal.vue`

**Step 1: Update AddSkillModal.vue**

Replace `web/src/components/AddSkillModal.vue` content:
```vue
<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { uploadSkill, installFromGithub } from '../api/skills'

const { t } = useI18n()

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

  const file = files[0]!
  if (!file.name.endsWith('.md')) {
    error.value = t('errors.onlyMdAllowed')
    return
  }

  loading.value = true
  error.value = ''

  try {
    await uploadSkill(file)
    emit('added')
    emit('close')
  } catch (e) {
    error.value = e instanceof Error ? e.message : t('errors.uploadFailed')
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
      error.value = t('errors.noSkillsInRepo')
    } else {
      emit('added')
      emit('close')
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : t('errors.installFailed')
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
        <h2 class="text-xl font-bold">{{ t('addModal.title') }}</h2>
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
          {{ t('addModal.uploadFile') }}
        </button>
        <button
          @click="tab = 'github'"
          :class="[
            'px-4 py-2 rounded-lg text-sm font-medium flex-1',
            tab === 'github' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700'
          ]"
        >
          {{ t('addModal.fromGithub') }}
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
          <p class="text-gray-600 mb-2">{{ t('addModal.dragDrop') }}</p>
          <p class="text-gray-400 text-sm mb-4">{{ t('addModal.or') }}</p>
          <label class="px-4 py-2 bg-blue-600 text-white rounded-lg cursor-pointer hover:bg-blue-700">
            {{ t('addModal.chooseFile') }}
            <input type="file" accept=".md" class="hidden" @change="onFileSelect" />
          </label>
        </div>
      </div>

      <div v-else>
        <input
          v-model="githubUrl"
          type="text"
          :placeholder="t('addModal.githubPlaceholder')"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <p class="mt-2 text-sm text-gray-500">
          {{ t('addModal.githubHelper') }}
        </p>
        <button
          @click="handleGithubInstall"
          :disabled="!githubUrl || loading"
          class="mt-4 w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
        >
          {{ t('addModal.install') }}
        </button>
      </div>

      <p v-if="error" class="mt-4 text-red-600 text-sm">{{ error }}</p>
      <p v-if="loading" class="mt-4 text-gray-600 text-sm">{{ t('status.loading') }}</p>
    </div>
  </div>
</template>
```

**Step 2: Commit component updates**

```bash
git add web/src/App.vue web/src/components/
git commit -m "feat: add i18n translations to all Vue components"
```

---

## Phase 3: Documentation

### Task 10: Build and verify localization works

**Step 1: Build the frontend**

```bash
cd /Users/wind/Desktop/projects/skill-router/web && npm run build
```

**Step 2: Start the backend server**

```bash
cd /Users/wind/Desktop/projects/skill-router && go run . &
```

**Step 3: Verify app loads at http://localhost:9527**

Use agent-browser to open http://localhost:9527 and verify:
- Page loads with English text
- Language toggle button visible in header
- Clicking toggle switches to Chinese

---

### Task 11: Capture screenshots

**Files:**
- Create: `docs/images/main-view.png`
- Create: `docs/images/add-skill-modal.png`

**Step 1: Create images directory**

```bash
mkdir -p /Users/wind/Desktop/projects/skill-router/docs/images
```

**Step 2: Use agent-browser to capture main view**

Navigate to http://localhost:9527 and take screenshot, save to `docs/images/main-view.png`

**Step 3: Use agent-browser to capture add modal**

Click "+ Add" button and take screenshot of modal, save to `docs/images/add-skill-modal.png`

---

### Task 12: Write English README

**Files:**
- Modify: `README.md`

**Step 1: Update README.md**

Replace `README.md` content:
```markdown
# Skill Router

A local web application to manage Claude Code skills and plugins.

![Main View](docs/images/main-view.png)

## Features

- **View all skills** from `~/.claude/commands/` and installed plugins
- **Enable/disable skills** individually or by plugin group
- **Delete skills** (user skills only)
- **Upload .md skill files** via drag-and-drop or file picker
- **Install skills from GitHub** repositories
- **Multi-language support** - English and Chinese with auto-detection

## Installation

Download the binary for your platform from [Releases](https://github.com/anthropics/skill-router/releases).

Or build from source:

```bash
make build
```

## Usage

```bash
./skill-router
```

This starts the server and opens your browser to http://localhost:9527

### Adding Skills

Click the **+ Add** button to:

1. **Upload a file** - Drag and drop or select a `.md` skill file
2. **Install from GitHub** - Enter a repository URL to install all skills from `.claude/commands/`

![Add Skill Modal](docs/images/add-skill-modal.png)

### Language

The interface automatically detects your browser language. Click the language toggle (EN/中) in the header to switch manually.

## Development

```bash
# Terminal 1: Run Go backend
go run .

# Terminal 2: Run Vue dev server
cd web && npm run dev
```

Then open http://localhost:5173 for hot-reload development.

### Project Structure

```
.
├── main.go                 # Entry point, HTTP server
├── internal/
│   ├── handler/            # HTTP handlers
│   ├── service/            # Business logic
│   └── config/             # Configuration management
├── web/                    # Vue 3 frontend
│   ├── src/
│   │   ├── components/     # Vue components
│   │   ├── i18n/           # Internationalization
│   │   └── api/            # API client
│   └── dist/               # Built frontend (embedded)
└── docs/
    └── images/             # Screenshots
```

## License

MIT
```

---

### Task 13: Write Chinese README

**Files:**
- Create: `README_CN.md`

**Step 1: Create README_CN.md**

Create `README_CN.md`:
```markdown
# Skill Router

一个用于管理 Claude Code 技能和插件的本地 Web 应用。

![主界面](docs/images/main-view.png)

## 功能特性

- **查看所有技能** - 来自 `~/.claude/commands/` 和已安装的插件
- **启用/禁用技能** - 支持单个技能或按插件组操作
- **删除技能** - 仅限用户技能
- **上传 .md 技能文件** - 支持拖放或文件选择
- **从 GitHub 安装技能** - 输入仓库地址即可安装
- **多语言支持** - 中文和英文，自动检测浏览器语言

## 安装

从 [Releases](https://github.com/anthropics/skill-router/releases) 下载适合您平台的二进制文件。

或从源码构建：

```bash
make build
```

## 使用方法

```bash
./skill-router
```

启动服务器后会自动打开浏览器访问 http://localhost:9527

### 添加技能

点击 **+ 添加** 按钮可以：

1. **上传文件** - 拖放或选择 `.md` 技能文件
2. **从 GitHub 安装** - 输入仓库 URL，自动安装 `.claude/commands/` 目录下的所有技能

![添加技能弹窗](docs/images/add-skill-modal.png)

### 语言切换

界面会自动检测您的浏览器语言。点击顶部的语言切换按钮（EN/中）可手动切换。

## 开发指南

```bash
# 终端 1：运行 Go 后端
go run .

# 终端 2：运行 Vue 开发服务器
cd web && npm run dev
```

然后访问 http://localhost:5173 进行热重载开发。

### 项目结构

```
.
├── main.go                 # 入口文件，HTTP 服务器
├── internal/
│   ├── handler/            # HTTP 处理器
│   ├── service/            # 业务逻辑
│   └── config/             # 配置管理
├── web/                    # Vue 3 前端
│   ├── src/
│   │   ├── components/     # Vue 组件
│   │   ├── i18n/           # 国际化
│   │   └── api/            # API 客户端
│   └── dist/               # 构建产物（已嵌入）
└── docs/
    └── images/             # 截图
```

## 许可证

MIT
```

**Step 2: Final commit**

```bash
git add README.md README_CN.md docs/images/
git commit -m "docs: add bilingual README with screenshots"
```

---

## Summary

**Total Tasks:** 13

**Phase 1 (Tasks 1-5):** i18n infrastructure - install vue-i18n, create translation files, wire to Vue app

**Phase 2 (Tasks 6-9):** Component updates - translate all UI text in App.vue, SkillCard.vue, PluginGroup.vue, AddSkillModal.vue

**Phase 3 (Tasks 10-13):** Documentation - build, verify, capture screenshots, write bilingual READMEs
