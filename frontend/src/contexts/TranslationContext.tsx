import { createContext, useContext, useMemo, type ReactNode } from 'react'
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
  const value = useMemo(() => ({ t, lang, setLanguage }), [t, lang, setLanguage])

  return (
    <TranslationContext.Provider value={value}>
      {children}
    </TranslationContext.Provider>
  )
}

export function useT() {
  return useContext(TranslationContext)
}
