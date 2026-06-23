import { useEffect, useState } from 'react'
import { Sun, Moon, Monitor } from 'lucide-react'
import { cn } from './cn'
import { useSettingsStore } from '../../stores/settingsStore'
import { useT } from '../../contexts/TranslationContext'

const themes = [
  { value: 'light', Icon: Sun },
  { value: 'dark', Icon: Moon },
  { value: 'system', Icon: Monitor },
] as const

export function ThemeToggle() {
  const { t } = useT()
  const { settings, saveSettings } = useSettingsStore()
  const currentTheme = settings?.theme ?? 'dark'
  const [open, setOpen] = useState(false)

  // Close dropdown on outside click
  useEffect(() => {
    if (!open) return
    const handler = () => setOpen(false)
    document.addEventListener('click', handler)
    return () => document.removeEventListener('click', handler)
  }, [open])

  const handleSelect = async (theme: string) => {
    await saveSettings({ ...settings!, theme })
    setOpen(false)
  }

  const activeTheme = themes.find((th) => th.value === currentTheme) ?? themes[1]
  const Icon = activeTheme.Icon

  return (
    <div className="relative">
      <button
        className={cn(
          'flex items-center gap-1.5 rounded-md px-2 py-1 text-xs',
          'text-neutral-500 hover:bg-neutral-100 dark:text-neutral-400 dark:hover:bg-neutral-800',
          'transition-colors',
        )}
        onClick={(e) => {
          e.stopPropagation()
          setOpen(!open)
        }}
        title={t(`theme.${activeTheme.value}`)}
      >
        <Icon className="h-3.5 w-3.5" />
        <span className="hidden sm:inline">{t(`theme.${activeTheme.value}`)}</span>
      </button>

      {open && (
        <div
          className={cn(
            'absolute right-0 top-full z-50 mt-1 min-w-[120px] overflow-hidden rounded-md border',
            'border-neutral-200 bg-white dark:border-neutral-700 dark:bg-neutral-900',
            'shadow-lg',
          )}
        >
          {themes.map(({ value, Icon: ItemIcon }) => (
            <button
              key={value}
              className={cn(
                'flex w-full items-center gap-2 px-3 py-1.5 text-xs',
                'hover:bg-neutral-100 dark:hover:bg-neutral-800',
                'transition-colors',
                value === currentTheme && 'text-blue-600 dark:text-blue-400',
              )}
              onClick={() => handleSelect(value)}
            >
              <ItemIcon className="h-3.5 w-3.5" />
              {t(`theme.${value}`)}
              {value === currentTheme && <span className="ml-auto">✓</span>}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
