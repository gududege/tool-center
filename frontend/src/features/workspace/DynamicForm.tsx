import { useState, useEffect, useCallback, useMemo, useRef, useLayoutEffect } from 'react'
import { createRoot, type Root } from 'react-dom/client'
import { withTheme } from '@rjsf/core'
import { Theme as ShadcnTheme, Templates as ShadcnTemplates } from '@rjsf/shadcn'
import { getDefaultFormState } from '@rjsf/utils'
import validator from '../../lib/rjsfValidator'
import { useT } from '../../contexts/TranslationContext'

import { Button } from '../../components/ui/Button'
import { ScrollArea } from '../../components/ui/ScrollArea'
import { openDirectoryDialog, openFileDialog, saveFileDialog } from '../../services/dialogApi'
import type { PluginDto, Tab, OutputEventDto, PluginHistoryEntryDto } from '../../types'
import { usePluginStore } from '../../stores/pluginStore'
import { useTabStore } from '../../stores/tabStore'
import { useOutputStore } from '../../stores/outputStore'
import { taskApi } from '../../services/taskApi'
import { pluginApi } from '../../services/pluginApi'
import { getShadcnDefaultSubThemeCss, appColorOverrideCss } from '../../vendor/rjsfShadcnThemes'

const DefaultObjectFieldTemplate = ShadcnTemplates.ObjectFieldTemplate!

// JSON Forms Group metadata carried into the RJSF uiSchema.
interface GroupInfo {
  label: string
  fields: string[]
}

// Themed shadcn Form component (only theme kept; antd/mantis/semantic-ui removed).
const ShadcnForm = withTheme(ShadcnTheme)

// ── Custom widgets for folder/file pickers ──
// RJSF widgets receive { value, onChange, schema, uiSchema, ... } — different shape from fields.
function FolderPickerWidget(props: any) {
  const { t } = useT()
  const { value, onChange } = props
  return (
    <div className="flex gap-2">
      <input
        type="text"
        className="flex-1 rounded-md border px-3 py-1.5 text-sm border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        value={value ?? ''}
        placeholder={t('form.select_directory')}
        onChange={(e) => onChange(e.target.value)}
      />
      <button
        type="button"
        className="rounded-md border px-2 py-1 text-sm border-input bg-background hover:bg-accent"
        onClick={async () => {
          const dir = await openDirectoryDialog({ title: t('form.select_directory') })
          if (dir) onChange(dir)
        }}
      >
        📁
      </button>
    </div>
  )
}

function FilePickerWidget(props: any) {
  const { t } = useT()
  const { value, onChange } = props
  return (
    <div className="flex gap-2">
      <input
        type="text"
        className="flex-1 rounded-md border px-3 py-1.5 text-sm border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        value={value ?? ''}
        placeholder={t('form.select_file')}
        onChange={(e) => onChange(e.target.value)}
      />
      <button
        type="button"
        className="rounded-md border px-2 py-1 text-sm border-input bg-background hover:bg-accent"
        onClick={async () => {
          const file = await openFileDialog({ title: t('form.select_file') })
          if (file) onChange(file)
        }}
      >
        📄
      </button>
    </div>
  )
}

function SaveFilePickerWidget(props: any) {
  const { t } = useT()
  const { value, onChange } = props
  return (
    <div className="flex gap-2">
      <input
        type="text"
        className="flex-1 rounded-md border px-3 py-1.5 text-sm border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        value={value ?? ''}
        placeholder={t('form.save_file')}
        onChange={(e) => onChange(e.target.value)}
      />
      <button
        type="button"
        className="rounded-md border px-2 py-1 text-sm border-input bg-background hover:bg-accent"
        onClick={async () => {
          const file = await saveFileDialog({ title: t('form.save_file') })
          if (file) onChange(file)
        }}
      >
        💾
      </button>
    </div>
  )
}

const customWidgets = {
  folderPicker: FolderPickerWidget,
  filePicker: FilePickerWidget,
  saveFilePicker: SaveFilePickerWidget,
}

// ── Custom object-field template with borders ──
//   * Root form with JSON Forms groups: each group renders as a bordered card.
//   * Nested objects: wrap the whole object in a subtle border.
//   * Root form without groups: stay borderless.
function BorderedObjectFieldTemplate(props: any) {
  const { fieldPathId, properties, uiSchema } = props
  const isRoot = fieldPathId?.$id === 'root'
  const groups: GroupInfo[] = uiSchema?.['ui:groups'] ?? []

  // Nested objects keep a subtle card-like border.
  if (!isRoot) {
    return (
      <div className="rounded-lg border border-border bg-card p-4">
        <DefaultObjectFieldTemplate {...props} />
      </div>
    )
  }

  // Root form: if the JSON Forms layout defines groups, render each group as a bordered card.
  if (groups.length > 0) {
    const groupedFieldNames = new Set(groups.flatMap((g) => g.fields))
    const remaining = properties.filter((p: any) => !groupedFieldNames.has(p.name))

    return (
      <div className="flex flex-col gap-4">
        {groups.map((group, idx) => (
          <div
            key={idx}
            className="rounded-lg border border-border bg-card p-4"
            >
            <h3 className="mb-3 text-sm font-medium text-foreground">
              {group.label}
            </h3>
            <div className="flex flex-col gap-4">
              {group.fields.map((fieldName) => {
                const element = properties.find((p: any) => p.name === fieldName)
                if (!element) return null
                return (
                  <div
                    key={element.name}
                    className={`${element.hidden ? 'hidden' : ''} flex`}
                  >
                    <div className="w-full">{element.content}</div>
                  </div>
                )
              })}
            </div>
          </div>
        ))}
        {remaining.length > 0 && (
          <div className="flex flex-col gap-4">
            {remaining.map((element: any) => (
              <div
                key={element.name}
                className={`${element.hidden ? 'hidden' : ''} flex`}
              >
                <div className="w-full">{element.content}</div>
              </div>
            ))}
          </div>
        )}
      </div>
    )
  }

  // Root form without groups: render with the default template (no outer border).
  return <DefaultObjectFieldTemplate {...props} />
}

// ── Shadow DOM host for the RJSF shadcn form ──
// The shadcn sub-theme CSS is Tailwind v4 (preflight + :root/.dark variable defs).
// Injecting it into the main document would override the app's --background/--foreground/
// --primary and recolor the sidebar/header. Rendering the form inside a shadow root keeps
// the Tailwind v4 CSS fully contained to the form area, while the app keeps its v3 vars.
//
// Theme unification: the form always uses the default sub-theme. An app-color override
// stylesheet (appColorOverrideCss) is injected AFTER the default CSS, replacing the
// shadow's shadcn oklch variables with the app's HSL values (from index.css) so the form
// area's colors match the rest of the app exactly — no separate form theme anymore.
//
// Dark mode: the CSS keys its dark variant off a `.dark` class, so we mirror
// <html>'s dark class onto the shadow host — flipping the app theme flips the form too.
interface ShadowFormHostProps {
  isDark: boolean
  schema: any
  uiSchema: any
  formData: Record<string, any>
  onChange: (formData: Record<string, any>) => void
  formRef: React.RefObject<any>
}

/** The actual form tree rendered into the shadow root. */
function ShadowFormContent({ schema, uiSchema, formData, onChange, formRef }: Omit<ShadowFormHostProps, 'isDark'>) {
  return (
      <ShadcnForm
        ref={formRef}
        schema={schema}
        uiSchema={uiSchema ?? {}}
        formData={formData}
        onChange={(e: any) => onChange(e.formData ?? {})}
        validator={validator}
        widgets={customWidgets}
        templates={{
          ObjectFieldTemplate: BorderedObjectFieldTemplate,
          ButtonTemplates: { SubmitButton: () => null },
        }}
        experimental_defaultFormStateBehavior={{
          arrayMinItems: { populate: 'never' },
          emptyObjectFields: 'skipDefaults',
        }}
        liveValidate={false}
        showErrorList={'top'}
        focusOnFirstError
        noHtml5Validate
      />
  )
}

function ShadowFormHost(props: ShadowFormHostProps) {
  const hostRef = useRef<HTMLDivElement | null>(null)
  const shadowRef = useRef<ShadowRoot | null>(null)
  const rootRef = useRef<Root | null>(null)
  const innerRef = useRef<HTMLDivElement | null>(null)

  // 1. Attach the shadow root once and inject the form stylesheets.
  useLayoutEffect(() => {
    if (!hostRef.current || shadowRef.current) return
    // StrictMode (dev) re-runs effects on the same host node; reuse an existing
    // shadow root instead of attachShadow() which throws "already exists".
    const existing = hostRef.current.shadowRoot
    const shadow = existing ?? hostRef.current.attachShadow({ mode: 'open' })
    shadowRef.current = shadow

    // Stylesheet 1: the default shadcn sub-theme. The CSS has been rewritten so its
    // :root{} → :host{} and .dark{} → :host(.dark){}, i.e. the shadcn variables are
    // defined on this shadow host (and inherited into the shadow tree). This provides
    // the Tailwind v4 preflight + utility classes the RJSF shadcn widgets rely on.
    const themeStyleEl = document.createElement('style')
    themeStyleEl.setAttribute('data-rjsf-shadcn-theme', 'default')
    themeStyleEl.textContent = getShadcnDefaultSubThemeCss()
    shadow.appendChild(themeStyleEl)

    // Stylesheet 2: app-color override. Injected AFTER the default CSS so it wins at
    // equal :host specificity. Replaces the sub-theme's oklch palette with the app's
    // HSL values (wrapped in hsl() — the v4 utilities consume var(--x) directly), so
    // the form area matches the rest of the app (neutral gray) instead of its own palette.
    const overrideStyleEl = document.createElement('style')
    overrideStyleEl.setAttribute('data-rjsf-app-override', 'true')
    overrideStyleEl.textContent = appColorOverrideCss
    shadow.appendChild(overrideStyleEl)

    // Inner wrapper — purely structural. It fills the host, inherits the app's
    // font, and uses --background/--foreground (inherited from :host) so the form
    // area is opaque in the active (light/dark) app palette.
    const inner = document.createElement('div')
    inner.setAttribute('data-rjsf-form-root', 'true')
    inner.style.height = '100%'
    inner.style.width = '100%'
    inner.style.color = 'var(--foreground)'
    inner.style.backgroundColor = 'var(--background)'
    inner.style.fontFamily = 'inherit'
    shadow.appendChild(inner)
    innerRef.current = inner

    // React root rendering the form inside the shadow tree.
    rootRef.current = createRoot(inner)
    rootRef.current.render(<ShadowFormContent {...props} />)

    return () => {
      rootRef.current?.unmount()
      rootRef.current = null
      shadowRef.current = null
      innerRef.current = null
    }
    // Attach once; subsequent prop changes handled by the effects below.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // 2. Mirror the app's dark mode onto the shadow host element so the
  //    :host(.dark) variable block flips the form in lockstep with the main UI.
  useEffect(() => {
    if (hostRef.current) {
      hostRef.current.classList.toggle('dark', props.isDark)
    }
  }, [props.isDark])

  // 3. Re-render the form tree when its props change.
  useEffect(() => {
    rootRef.current?.render(<ShadowFormContent {...props} />)
  })

  return <div ref={hostRef} className="h-full w-full" />
}

// 根据插件参数映射与当前表单数据生成命令行参数字符串。
function buildCommand(plugin: PluginDto | null, formData: Record<string, any>): string {
  if (!plugin) return ''
  const exe = plugin.execution?.exe ?? ''
  const params = plugin.execution?.parameters ?? []
  const args: string[] = []

  for (const p of params) {
    const val = formData[p.field]
    if (val === undefined || val === null || val === '') continue

    switch (p.kind) {
      case 'argument':
        args.push(String(val))
        break
      case 'argument-array':
        if (Array.isArray(val)) {
          for (const v of val) args.push(String(v))
        }
        break
      case 'option':
        args.push(p.flag ?? `--${p.field}`, String(val))
        break
      case 'option-array': {
        if (!Array.isArray(val)) break
        if (p.style === 'join' || p.style === 'equals') {
          const sep = p.separator || ','
          const joined = val.map(String).join(sep)
          if (p.style === 'equals') {
            args.push(`${p.flag}=${joined}`)
          } else {
            args.push(p.flag ?? `--${p.field}`, joined)
          }
        } else {
          for (const v of val) {
            args.push(p.flag ?? `--${p.field}`, String(v))
          }
        }
        break
      }
      case 'switch':
        if (val === true || val === 'true' || val === '1') {
          args.push(p.flag ?? `--${p.field}`)
        }
        break
      case 'bool-option':
        args.push(p.flag ?? `--${p.field}`, val === true || val === 'true' || val === '1' ? 'true' : 'false')
        break
      case 'dual-switch':
        if (val === true || val === 'true' || val === '1') {
          args.push(p.trueFlag ?? `--${p.field}`)
        } else if (p.falseFlag) {
          args.push(p.falseFlag)
        }
        break
    }
  }

  return `${exe} ${args.join(' ')}`.trim()
}

interface DynamicFormProps {
  tab: Tab
}

export function DynamicForm({ tab }: DynamicFormProps) {
  const { t, lang } = useT()
  const { fetchPlugin } = usePluginStore()
  const { setTabTaskId, updateTabTitle } = useTabStore()
  const { appendOutput } = useOutputStore()

  const [plugin, setPlugin] = useState<PluginDto | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [formData, setFormData] = useState<Record<string, any>>({})
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

  // 开启 preview 后实时计算命令行参数字符串。
  const previewCommand = useMemo(() => buildCommand(plugin, formData), [plugin, formData])

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
  // RJSF is configured with emptyObjectFields: 'skipDefaults' so clearing a
  // field does not re-inject schema defaults. We manually compute and set the
  // initial defaults once when the plugin is first loaded.
  useEffect(() => {
    if (!plugin) return
    if (defaultsAppliedFor.current === plugin.metadata.id) return
    const defaults = getDefaultFormState(validator, localizedSchema, {}, localizedSchema)
    setFormData((prev) => ({ ...(defaults as Record<string, any> || {}), ...prev }))
    defaultsAppliedFor.current = plugin.metadata.id
  }, [plugin, localizedSchema])

  // ── Update tab title on language change ──
  useEffect(() => {
    if (!plugin) return
    const localizedName = lang === 'zh' && plugin.metadata.name_cn ? plugin.metadata.name_cn : plugin.metadata.name
    updateTabTitle(tab.id, localizedName)
  }, [plugin, lang, tab.id, updateTabTitle])

  // ── Load historical output ──
  useEffect(() => {
    if (tab.taskId) {
      taskApi.getTaskOutput(tab.taskId).then((lines: OutputEventDto[]) => {
        for (const line of lines) {
          appendOutput(tab.taskId!, line)
        }
      })
    }
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
      await pluginApi.saveHistory(plugin.metadata.id, formData)
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
        formData,
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

  return (
    <div className="flex h-full flex-col">
      {/* Form title */}
      <div className="border-b border-border px-4 py-3">
        <h2 className="text-base font-semibold">
          {lang === 'zh' && plugin.metadata.name_cn ? plugin.metadata.name_cn : plugin.metadata.name}
        </h2>
        {(() => {
          const d = lang === 'zh' && plugin.metadata.description_cn ? plugin.metadata.description_cn : plugin.metadata.description
          return d ? <p className="mt-0.5 text-xs text-muted-foreground">{d}</p> : null
        })()}
      </div>

      {/* Form area — RJSF shadcn form rendered in a Shadow DOM so its Tailwind v4
          CSS stays contained and never recolors the rest of the app. */}
      <ScrollArea className="flex-1 p-4">
        <div className="mx-auto max-w-2xl">
          <ShadowFormHost
            isDark={isDark}
            schema={localizedSchema}
            uiSchema={localizedUiSchema}
            formData={formData}
            onChange={(data) => setFormData(data)}
            formRef={formRef}
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
        <Button onClick={handleRun} disabled={running}>
          {running ? t('form.running') : t('form.run')}
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
                      key={idx}
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

// ── Localization helpers ──

function localizeSchema(schema: any, lang: string): any {
  if (!schema || !schema.properties) return schema
  const out = { ...schema }
  if (schema.title_cn && lang === 'zh') out.title = schema.title_cn
  if (schema.description_cn && lang === 'zh') out.description = schema.description_cn

  const props: Record<string, any> = {}
  for (const [key, val] of Object.entries(schema.properties)) {
    const prop = { ...(val as any) }
    if (lang === 'zh' && prop.title_cn) prop.title = prop.title_cn
    if (lang === 'zh' && prop.description_cn) prop.description = prop.description_cn
    // Localize anyOf/oneOf titles
    if (prop.anyOf) {
      prop.anyOf = prop.anyOf.map((v: any) => {
        const item = { ...v }
        if (lang === 'zh' && v.title_cn) item.title = v.title_cn
        return item
      })
    }
    if (prop.oneOf) {
      prop.oneOf = prop.oneOf.map((v: any) => {
        const item = { ...v }
        if (lang === 'zh' && v.title_cn) item.title = v.title_cn
        return item
      })
    }
    // Localize enum labels — for now keep values as-is, they're identifiers
    props[key] = prop
  }
  out.properties = props
  return out
}

// ── uiSchema 转换 ──
//
// 插件的 uischema.json 用的是 JSON Forms 的布局格式（type: VerticalLayout / Group /
// Control，scope: "#/properties/<field>"），但 RJSF 用的是以字段名为 key 的对象式
// uiSchema（{ indir: { 'ui:widget': 'folderPicker' } }）。RJSF 不解析 elements/scope，
// 所以必须在这里把 JSON Forms 布局展平成 RJSF uiSchema，否则 options.folderPicker /
// saveFilePicker / inputType 这些自定义 widget 配置全部失效，字段退回 RJSF 默认
// widget（例如 format: data-url 会变成文件上传框，而非目录选择框）。
//
// 展平同时收集 Control 的出现顺序，写入 ui:order，强制 RJSF 按布局顺序渲染字段
// （RJSF 默认按 schema properties 顺序，但布局可能重排，ui:order 保证一致）。
// JSON Forms 的 Group 分组被保留到 'ui:groups' 中，自定义 ObjectFieldTemplate 据此
// 把字段渲染成带边框的卡片。

function localizeUiSchema(uiSchema: any, lang: string): any {
  if (!uiSchema) return undefined
  const fieldUi: Record<string, any> = {}
  const order: string[] = []
  const groups: GroupInfo[] = []
  walkJsonFormsLayout(uiSchema, fieldUi, order, groups, lang)
  if (Object.keys(fieldUi).length === 0 && groups.length === 0) return undefined
  // ui:order 末尾的 '*' 表示「其余字段按 schema 顺序追加」，避免漏掉未在布局里
  // 显式出现的字段。
  const out: Record<string, any> = { ...fieldUi, 'ui:order': [...order, '*'] }
  if (groups.length > 0) {
    out['ui:groups'] = groups
  }
  return out
}

function walkJsonFormsLayout(
  el: any,
  fieldUi: Record<string, any>,
  order: string[],
  groups: GroupInfo[],
  lang: string,
) {
  if (!el || typeof el !== 'object') return
  const type: string = el.type

  // Group —— 收集字段并保留分组信息，然后递归处理子 Control。
  if (type === 'Group') {
    const groupFields: string[] = []
    if (Array.isArray(el.elements)) {
      for (const child of el.elements) {
        const field = collectControlField(child, fieldUi, order, lang)
        if (field && !groupFields.includes(field)) groupFields.push(field)
      }
    }
    if (groupFields.length > 0) {
      const label = lang === 'zh' && el.label_cn ? el.label_cn : el.label
      groups.push({ label: label || '', fields: groupFields })
    }
    return
  }

  // 布局容器：VerticalLayout / HorizontalLayout / Categorization / Category。
  if (type === 'VerticalLayout' || type === 'HorizontalLayout' || type === 'Categorization' || type === 'Category') {
    if (Array.isArray(el.elements)) {
      for (const child of el.elements) walkJsonFormsLayout(child, fieldUi, order, groups, lang)
    }
    return
  }

  // Control —— 对应一个字段。
  collectControlField(el, fieldUi, order, lang)

  // 兜底：若顶层就是以字段名为 key 的 RJSF uiSchema（非 JSON Forms 布局），直接透传。
  // 这种情况没有 type 字段，fieldUi 保持为空，groups 也为空，调用方会原样返回 uiSchema。
}

function collectControlField(
  el: any,
  fieldUi: Record<string, any>,
  order: string[],
  lang: string,
): string | null {
  if (!el || typeof el !== 'object') return null
  const type: string = el.type
  if (type !== 'Control') return null
  const field = scopeToField(el.scope)
  if (!field) return null
  const ui: any = {}
  const opts = el.options
  if (opts && typeof opts === 'object') {
    if (opts.folderPicker === true) {
      ui['ui:widget'] = 'folderPicker'
    } else if (opts.saveFilePicker === true) {
      ui['ui:widget'] = 'saveFilePicker'
    } else if (opts.filePicker === true) {
      ui['ui:widget'] = 'filePicker'
    } else if (typeof opts.widget === 'string' && opts.widget) {
      // 通用 RJSF widget 映射：textarea / radio / select / range / hidden 等。
      ui['ui:widget'] = opts.widget
    }
    if (opts.inputType) {
      // RJSF TextWidget 读 ui:inputType 决定 <input type>（password/number/...）。
      ui['ui:inputType'] = opts.inputType
    }
  }
  fieldUi[field] = ui
  if (!order.includes(field)) order.push(field)
  return field
}

// scopeToField 把 JSON Forms 的 scope "#/properties/indir" 提取出字段名 "indir"。
// 支持嵌套 "#/properties/foo/properties/bar" → "foo"（RJSF 嵌套对象用点号访问，
// 这里只取顶层字段名，足够当前插件 schema 的扁平结构）。
function scopeToField(scope: any): string | null {
  if (typeof scope !== 'string') return null
  // 形如 #/properties/<name>
  const m = scope.match(/^#\/properties\/([^/]+)/)
  return m ? m[1] : null
}
