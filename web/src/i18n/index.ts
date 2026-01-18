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
