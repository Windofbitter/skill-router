<script setup lang="ts">
import type { Skill } from '../types/skill'

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

const isUserSkill = props.skill.source === 'user'
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
          {{ skill.description || 'No description' }}
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
        {{ skill.enabled ? 'Enabled' : 'Disabled' }}
      </span>
    </div>

    <!-- Only show action buttons for user skills -->
    <div v-if="isUserSkill" class="mt-4 flex gap-2">
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

    <!-- Plugin skills can be toggled but not deleted -->
    <div v-else class="mt-4 flex items-center gap-2">
      <button
        v-if="skill.enabled"
        @click="emit('disablePlugin', skill.pluginName, skill.fileName)"
        class="px-3 py-1 text-sm bg-yellow-100 text-yellow-800 rounded hover:bg-yellow-200"
      >
        Disable
      </button>
      <button
        v-else
        @click="emit('enablePlugin', skill.pluginName, skill.fileName)"
        class="px-3 py-1 text-sm bg-green-100 text-green-800 rounded hover:bg-green-200"
      >
        Enable
      </button>
      <span class="text-xs text-gray-400">from {{ skill.pluginName }}</span>
    </div>
  </div>
</template>
