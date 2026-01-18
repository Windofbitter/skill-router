export interface Skill {
  name: string
  description: string
  fileName: string
  filePath: string
  enabled: boolean
  source: 'user' | 'plugin'
  pluginName: string
}
