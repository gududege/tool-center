import { useState, useEffect, useCallback, useMemo, useRef } from 'react'
import { getDefaultFormState } from '@rjsf/utils'
import validator from '../../lib/rjsfValidator'
import { useT } from '../../contexts/TranslationContext'
import { Button } from '../../components/ui/Button'
import { ScrollArea } from '../../components/ui/ScrollArea'
import type { PluginDto, Tab, OutputEventDto, PluginHistoryEntryDto } from '../../types'
import { usePluginStore } from '../../stores/pluginStore'
import { useTabStore } from '../../stores/tabStore'
import { useOutputStore } from '../../stores/outputStore'
import { useSettingsStore } from '../../stores/settingsStore'
import { taskApi } from '../../services/taskApi'
import { pluginApi } from '../../services/pluginApi'
import { ShadowFormHost, BorderedObjectFieldTemplate } from '../../components/ui/ShadowFormHost'
import { FolderPickerWidget } from '../../components/widgets/FolderPickerWidget'
import { FilePickerWidget } from '../../components/widgets/FilePickerWidget'
import { SaveFilePickerWidget } from '../../components/widgets/SaveFilePickerWidget'
import { buildCommand } from '../../lib/buildCommand'
import { localizeSchema, localizeUiSchema } from '../../lib/pluginLocalization'

const customWidgets = {
  folderPicker: FolderPickerWidget,
  filePicker: FilePickerWidget,
  saveFilePicker: SaveFilePickerWidget,
}

interface DynamicFormProps {
  tab: Tab
}

export function DynamicForm({ tab }: DynamicFormProps) {
  const { t, lang } = useT()
  const { fetchPlugin } = usePluginStore()
  const { setTabTaskId, updateTabTitle } = useTabStore()
  const { appendOutput } = useOutputStore()
  const settings = useSettingsStore((s) => s.settings)

  const [plugin, setPlugin] = useState<PluginDto | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [formData, setFormData] = useState<Record<string, unknown>>({})
  const [running, setRunning] = useState(false)
  const [previewEnabled, setPreviewEnabled] = useState(false)
  const [history, setHistory] = useState<PluginHistoryEntryDto[]>([])
  const [historyOpen, setHistoryOpen] = useState(false)

  // RJSF Form 实例引用 —— 用于在点 Run 时程序化触发校验。
  // 因为表单渲染在 Shadow DOM 里，外层 Run 按钮无法依赖 RJSF 默认的 SubmitButton
  // 触发校验，所以持有实例后手动调用 validateForm()。
  const formRef = useRef<any>(null)

  // 历史下拉容器引用，用于点击外部关闭。
  const historyRef = useRef<HTMLDivElement | null>(null)

  // Track which plugin already had its schema defaults applied so we only fill
  // defaults once per plugin/tab instance. After this, skipDefaults prevents
  // RJSF from re-injecting defaults when the user clears a field.
  const defaultsAppliedFor = useRef<string | null>(null)

  // Track the app's dark mode so the shadow form can mirror it.
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

  // 点击历史下拉外部时关闭下拉菜单。
  useEffect(() => {
    if (!historyOpen) return
    const handleClick = (e: MouseEvent) => {
      if (!historyRef.current?.contains(e.target as Node)) {
        setHistoryOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [historyOpen])

  // ── Localize schema / uiSchema ──
  const localizedSchema = useMemo(() => {
    if (!plugin) return {}
    const schema = plugin.form?.schema ?? {}
    return localizeSchema(schema, lang)
  }, [plugin, lang])

  const localizedUiSchema = useMemo(() => {
    if (!plugin) return undefined
    return localizeUiSchema(plugin.form?.uiSchema, lang)
  }, [plugin, lang])

  // ── Load plugin ──
  useEffect(() => {
    let cancelled = false
    setLoading(true)
    fetchPlugin(tab.pluginId)
      .then((p) => {
        if (!cancelled) {
          setPlugin(p)
          setLoading(false)
        }
      })
      .catch((err) => {
        if (!cancelled) {
          setError(err?.message ?? t('form.error_load'))
          setLoading(false)
        }
      })
    return () => { cancelled = true }
  }, [tab.pluginId])

  // ── Apply schema defaults once per plugin instance ──
  const resolvedFormData = useMemo(() => {
    if (!plugin || !localizedSchema) return formData
    if (defaultsAppliedFor.current === plugin.metadata.id) return formData
    const defaults = getDefaultFormState(validator, localizedSchema, {}, localizedSchema)
    defaultsAppliedFor.current = plugin.metadata.id
    return { ...(defaults as Record<string, unknown> || {}), ...formData }
  }, [plugin, localizedSchema, formData])

  const runnerPaths = settings?.runnerPaths
  const runnerOptions = settings?.runnerOptions
  const previewCommand = useMemo(
    () => buildCommand(plugin, resolvedFormData, runnerPaths, runnerOptions),
    [plugin, resolvedFormData, runnerPaths, runnerOptions],
  )

  // ── Update tab title on language change ──
  useEffect(() => {
    if (!plugin) return
    const localizedName = lang === 'zh' && plugin.metadata.name_cn ? plugin.metadata.name_cn : plugin.metadata.name
    updateTabTitle(tab.id, localizedName)
  }, [plugin, lang, tab.id, updateTabTitle])

  // ── Load historical output ──
  useEffect(() => {
    let cancelled = false
    if (tab.taskId) {
      taskApi.getTaskOutput(tab.taskId).then((lines: OutputEventDto[]) => {
        if (cancelled) return
        for (const line of lines) {
          appendOutput(tab.taskId!, line)
        }
      })
    }
    return () => { cancelled = true }
  }, [tab.taskId])

  // ── Load plugin history ──
  useEffect(() => {
    if (!plugin) return
    let cancelled = false
    pluginApi.getHistory(plugin.metadata.id).then((entries) => {
      if (!cancelled) setHistory(entries)
    })
    return () => { cancelled = true }
  }, [plugin])

  // ── Run plugin ──
  const handleRun = useCallback(async () => {
    if (!plugin) return
    // 提交前先程序化校验表单。validateForm() 返回 false 时 RJSF 会把错误
    // 写入内部 state 并通过 showErrorList 渲染（focusOnFirstError 聚焦首个错误字段）。
    // 校验不通过则不保存历史、不发起后端运行。
    if (formRef.current && typeof formRef.current.validateForm === 'function') {
      const ok = formRef.current.validateForm()
      if (!ok) return
    }

    // 校验通过后保存当前参数到插件目录 history.json（最多 5 条滚动保存）。
    try {
      await pluginApi.saveHistory(plugin.metadata.id, resolvedFormData)
      // 保存成功后刷新本地历史列表
      const updated = await pluginApi.getHistory(plugin.metadata.id)
      setHistory(updated)
    } catch (err: any) {
      // 历史保存失败不影响运行；可在控制台输出警告。
      // eslint-disable-next-line no-console
      console.warn('save history failed', err)
    }

    setRunning(true)
    try {
      const resp = await taskApi.runPlugin({
        pluginId: plugin.metadata.id,
        formData: resolvedFormData,
      })
      setTabTaskId(tab.id, resp.taskId)
    } catch (err: any) {
      setError(err?.message ?? t('form.error_run'))
    } finally {
      setRunning(false)
    }
  }, [plugin, formData, tab.id])

  const handleReset = useCallback(() => {
    setFormData({})
    setPreviewEnabled(false)
  }, [])

  // ── Command preview toggle ──
  const handleTogglePreview = useCallback(() => {
    setPreviewEnabled((prev) => !prev)
  }, [])

  // ── Import / Export ──
  const handleExport = useCallback(() => {
    const json = JSON.stringify(formData, null, 2)
    const blob = new Blob([json], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${plugin?.metadata.id ?? 'params'}-${Date.now()}.json`
    a.click()
    URL.revokeObjectURL(url)
  }, [formData, plugin])

  const handleImport = useCallback(() => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = '.json'
    input.onchange = () => {
      const file = input.files?.[0]
      if (!file) return
      const reader = new FileReader()
      reader.onload = () => {
        try {
          const data = JSON.parse(reader.result as string)
          setFormData((prev) => ({ ...prev, ...data }))
        } catch {
          // ignore invalid JSON
        }
      }
      reader.readAsText(file)
    }
    input.click()
  }, [])

  const handleHistoryLoad = useCallback((entry: PluginHistoryEntryDto) => {
    setFormData(entry.formData)
    setHistoryOpen(false)
  }, [])

  if (loading) {
    return (
      <div className="flex h-full items-center justify-center text-sm text-neutral-400">
        {t('form.loading')}
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex h-full items-center justify-center text-sm text-red-500">
        {error}
      </div>
    )
  }

  if (!plugin) {
    return (
      <div className="flex h-full items-center justify-center text-sm text-neutral-400">
        {t('form.not_found')}
      </div>
    )
  }

  const localizedDescription = lang === 'zh' && plugin.metadata.description_cn
    ? plugin.metadata.description_cn
    : plugin.metadata.description

  return (
    <div className="flex h-full flex-col">
      {/* Form title */}
      <div className="border-b border-border px-4 py-3">
        <h2 className="text-base font-semibold">
          {lang === 'zh' && plugin.metadata.name_cn ? plugin.metadata.name_cn : plugin.metadata.name}
        </h2>
        {localizedDescription && (
          <p className="mt-0.5 text-xs text-muted-foreground">{localizedDescription}</p>
        )}
      </div>

      {/* Form area — RJSF shadcn form rendered in a Shadow DOM so its Tailwind v4
          CSS stays contained and never recolors the rest of the app. */}
      <ScrollArea className="flex-1 p-4">
        <div className="mx-auto max-w-2xl">
          <ShadowFormHost
            isDark={isDark}
            schema={localizedSchema}
            uiSchema={localizedUiSchema}
            formData={resolvedFormData}
            onChange={(data) => setFormData(data)}
            formRef={formRef}
            widgets={customWidgets}
            templates={{
              ObjectFieldTemplate: BorderedObjectFieldTemplate,
              ButtonTemplates: { SubmitButton: () => null },
            }}
          />

          {/* Command preview */}
          {previewEnabled && (
            <div className="mt-4 rounded-md bg-neutral-100 p-3 font-mono text-xs dark:bg-neutral-800">
              <div className="mb-1 text-muted-foreground">{t('form.preview_label')}</div>
              <code>{previewCommand}</code>
            </div>
          )}
        </div>
      </ScrollArea>

      {/* Action buttons */}
      <div className="flex items-center gap-2 border-t border-border px-4 py-3">
        <Button onClick={handleRun} disabled={running} className="min-w-24">
          {t('form.run')}
        </Button>
        <Button variant="secondary" onClick={handleReset}>
          {t('form.reset')}
        </Button>
        <Button
          variant={previewEnabled ? 'secondary' : 'ghost'}
          onClick={handleTogglePreview}
          title={previewEnabled ? t('form.preview_on') : t('form.preview')}
        >
          {previewEnabled ? t('form.preview_on') : t('form.preview')}
        </Button>
        <div className="mx-1 h-5 w-px bg-border" />
        <Button variant="outline" onClick={handleImport} title={t('form.import')}>
          {t('form.import')}
        </Button>
        <Button variant="outline" onClick={handleExport} title={t('form.export')}>
          {t('form.export')}
        </Button>
        <div className="relative" ref={historyRef}>
          <Button
            variant="outline"
            onClick={() => setHistoryOpen((prev) => !prev)}
            title={t('form.history')}
          >
            {t('form.history')}
          </Button>
          {historyOpen && (
            <div className="absolute right-0 bottom-full z-50 mb-1 w-56 rounded-md border border-border bg-background shadow-md">
              <div className="max-h-60 overflow-auto py-1">
                {history.length === 0 ? (
                  <div className="px-3 py-2 text-xs text-muted-foreground">
                    {t('form.no_history')}
                  </div>
                ) : (
                  history.map((entry, idx) => (
                    <button
                      key={entry.timestamp}
                      type="button"
                      className="w-full px-3 py-2 text-left text-xs hover:bg-accent"
                      onClick={() => handleHistoryLoad(entry)}
                      title={entry.label}
                    >
                      {entry.label}
                    </button>
                  ))
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

// ── Localization helpers ── imported from '../../lib/pluginLocalization'
