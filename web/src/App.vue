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
