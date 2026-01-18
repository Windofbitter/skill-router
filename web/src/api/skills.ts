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
