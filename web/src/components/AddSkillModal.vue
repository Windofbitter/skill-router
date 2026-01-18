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

  const file = files[0]!
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
