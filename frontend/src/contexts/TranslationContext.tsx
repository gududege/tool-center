import { createContext, useContext, type ReactNode } from 'react'
import { useTranslation } from '../hooks/useTranslation'

interface TranslationContextValue {
  t: (key: string, fallback?: string) => string
  lang: string
  setLanguage: (lang: string) => void
}

const TranslationContext = createContext<TranslationContextValue>({
  t: (k) => k,
  lang: 'en',
  setLanguage: () => {},
})

export function TranslationProvider({ children }: { children: ReactNode }) {
  const { t, lang, setLanguage } = useTranslation()

  return (
    <TranslationContext.Provider value={{ t, lang, setLanguage }}>
      {children}
    </TranslationContext.Provider>
  )
}

export function useT() {
  return useContext(TranslationContext)
}
