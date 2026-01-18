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
