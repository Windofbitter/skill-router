<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Skill } from '../types/skill'
import SkillCard from './SkillCard.vue'

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
            {{ enabledCount }}/{{ skills.length }} skills enabled
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          v-if="!allEnabled"
          @click="emit('enablePlugin', pluginName)"
          class="px-3 py-1 text-sm bg-green-100 text-green-800 rounded hover:bg-green-200"
        >
          Enable All
        </button>
        <button
          v-if="!allDisabled"
          @click="emit('disablePlugin', pluginName)"
          class="px-3 py-1 text-sm bg-yellow-100 text-yellow-800 rounded hover:bg-yellow-200"
        >
          Disable All
        </button>
        <button
          @click="emit('deletePlugin', pluginName)"
          class="px-3 py-1 text-sm bg-red-100 text-red-800 rounded hover:bg-red-200"
        >
          Remove
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
