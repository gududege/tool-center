import { useState, useEffect, useCallback, useRef, useMemo } from 'react'
import { ShadowFormHost } from '../../components/ui/ShadowFormHost'
import { useSettingsStore } from '../../stores/settingsStore'
import { useT } from '../../contexts/TranslationContext'
import { getSettingsSchema, settingsUiSchema } from './settingsSchema'
import { FolderPickerWidget } from '../../components/widgets/FolderPickerWidget'
import { FilePickerWidget } from '../../components/widgets/FilePickerWidget'

const settingWidgets = {
  folderPicker: FolderPickerWidget,
  filePicker: FilePickerWidget,
}

interface SettingsPageProps {
  open: boolean
  onClose: () => void
}

export function SettingsPage({ open, onClose }: SettingsPageProps) {
  const { settings, saveSettings } = useSettingsStore()
  const formRef = useRef<any>(null)
  const { t } = useT()

  const [isDark, setIsDark] = useState(
    typeof document !== 'undefined' && document.documentElement.classList.contains('dark'),
  )
  useEffect(() => {
    const el = document.documentElement
    const sync = () => setIsDark(el.classList.contains('dark'))
    sync()
    const observer = new MutationObserver(sync)
    observer.observe(el, { attributes: true, attributeFilter: ['class'] })
    return () => observer.disconnect()
  }, [])

  const formData = settings
    ? {
        theme: settings.theme,
        formTheme: settings.formTheme ?? 'default',
        language: settings.language ?? 'en',
        pluginDirectory: settings.pluginDirectory,
        runnerPaths: settings.runnerPaths ?? {},
        runnerOptions: settings.runnerOptions ?? {},
      }
    : {}

  const handleSubmit = useCallback(
    async (data: Record<string, unknown>) => {
      if (!settings) return
      await saveSettings({
        ...settings,
        theme: data.theme as string,
        formTheme: data.formTheme as string,
        language: data.language as string,
        pluginDirectory: data.pluginDirectory as string,
        runnerPaths: (data.runnerPaths as Record<string, string>) ?? {},
        runnerOptions: (data.runnerOptions as Record<string, string>) ?? {},
      })
      onClose()
    },
    [settings, saveSettings, onClose],
  )

  const DEFAULT_FORM_VALUES = {
    theme: 'dark',
    formTheme: 'default',
    language: 'en',
    pluginDirectory: './plugins',
    runnerPaths: {
      cmd: 'cmd',
      bat: 'cmd',
      powershell: 'powershell',
      sh: 'bash',
      python: 'python',
      nodejs: 'node',
    },
    runnerOptions: {
      cmd: '/C',
      bat: '/C',
      powershell: '-ExecutionPolicy Bypass -File',
    },
  }

  const handleResetDefaults = useCallback(async () => {
    if (!window.confirm(t('button.reset_settings_confirm'))) return
    if (!settings) return
    await saveSettings({
      ...settings,
      ...DEFAULT_FORM_VALUES,
    })
    onClose()
  }, [settings, saveSettings, onClose, t])

  const schema = useMemo(() => getSettingsSchema((key: string) => t(key)), [t])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="relative w-full max-w-xl max-h-[85vh] flex flex-col rounded-lg border border-border bg-background shadow-xl">
        <div className="flex items-center justify-between border-b border-border px-4 py-3">
          <h2 className="text-sm font-semibold">{t('settings.title')}</h2>
          <button
            onClick={onClose}
            className="rounded-md px-2 py-1 text-xs text-muted-foreground hover:bg-accent"
            title={t('button.close_settings')}
          >
            {t('settings.close')}
          </button>
        </div>
        <div className="flex-1 overflow-auto p-4">
          <ShadowFormHost
            isDark={isDark}
            schema={schema}
            uiSchema={settingsUiSchema}
            formData={formData}
            onChange={() => {}}
            formRef={formRef}
            widgets={settingWidgets}
            onSubmit={handleSubmit}
          />
          <div className="mt-4 flex items-center justify-between gap-2 border-t border-border pt-4">
            <button
              onClick={handleResetDefaults}
              className="rounded-md border border-destructive/30 px-3 py-1.5 text-xs text-destructive hover:bg-destructive/5 transition-colors"
            >
              {t('button.reset_settings')}
            </button>
            <button
              onClick={() => formRef.current?.submit()}
              className="rounded-md bg-primary px-4 py-1.5 text-xs text-primary-foreground hover:bg-primary/90 transition-colors"
            >
              {t('button.save_settings')}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
