/**
 * Simple i18n hook — no external dependencies.
 * Defaults to English, loads from JSON files.
 */

import { useCallback, useEffect, useState } from 'react'
import en from '../locales/en.json'
import zh from '../locales/zh.json'

type LocaleMessages = Record<string, string>

const locales: Record<string, LocaleMessages> = { en, zh }

function loadMessages(lang: string): LocaleMessages {
  return locales[lang] ?? en
}

function getInitialLanguage(): string {
  try {
    return localStorage.getItem('language') || 'en'
  } catch {
    return 'en'
  }
}

export function useTranslation() {
  const [lang, setLangState] = useState(getInitialLanguage)
  const [messages, setMessages] = useState<LocaleMessages>(() => loadMessages(getInitialLanguage()))

  useEffect(() => {
    const m = loadMessages(lang)
    setMessages(m)
    try {
      localStorage.setItem('language', lang)
    } catch { /* ignore */ }
  }, [lang])

  const setLanguage = useCallback((l: string) => {
    setLangState(l)
    // Dispatch custom event so other components can react
    window.dispatchEvent(new CustomEvent('languagechange', { detail: l }))
  }, [])

  const t = useCallback(
    (key: string, fallback?: string): string => {
      return messages[key] ?? fallback ?? key
    },
    [messages],
  )

  return { t, lang, setLanguage }
}
