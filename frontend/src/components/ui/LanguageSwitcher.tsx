import { useT } from '../../contexts/TranslationContext'
import { cn } from './cn'

export function LanguageSwitcher() {
  const { t, lang, setLanguage } = useT()

  return (
    <div className="flex items-center gap-0.5 rounded-md border border-border bg-background text-xs">
      <button
        className={cn(
          'rounded-l-md px-2 py-1 transition-colors',
          lang === 'en'
            ? 'bg-primary text-primary-foreground'
            : 'text-muted-foreground hover:bg-neutral-100 dark:hover:bg-neutral-800',
        )}
        onClick={() => setLanguage('en')}
        title={t('common.lang.en')}
      >
        EN
      </button>
      <button
        className={cn(
          'rounded-r-md px-2 py-1 transition-colors',
          lang === 'zh'
            ? 'bg-primary text-primary-foreground'
            : 'text-muted-foreground hover:bg-neutral-100 dark:hover:bg-neutral-800',
        )}
        onClick={() => setLanguage('zh')}
        title={t('common.lang.zh')}
      >
        中
      </button>
    </div>
  )
}
