import { useEffect, useLayoutEffect, useRef } from 'react'
import { createRoot, type Root } from 'react-dom/client'
import { withTheme } from '@rjsf/core'
import { Theme as ShadcnTheme, Templates as ShadcnTemplates } from '@rjsf/shadcn'
import validator from '../../lib/rjsfValidator'
import { getShadcnDefaultSubThemeCss, appColorOverrideCss } from '../../vendor/rjsfShadcnThemes'

const DefaultObjectFieldTemplate = ShadcnTemplates.ObjectFieldTemplate!

const ShadcnForm = withTheme(ShadcnTheme)

export interface GroupInfo {
  label: string
  fields: string[]
}

export function BorderedObjectFieldTemplate(props: any) {
  const { fieldPathId, properties, uiSchema } = props
  const isRoot = fieldPathId?.$id === 'root'
  const groups: GroupInfo[] = uiSchema?.['ui:groups'] ?? []

  if (!isRoot) {
    return (
      <div className="rounded-lg border border-border bg-card p-4">
        <DefaultObjectFieldTemplate {...props} />
      </div>
    )
  }

  if (groups.length > 0) {
    const groupedFieldNames = new Set(groups.flatMap((g) => g.fields))
    const remaining = properties.filter((p: any) => !groupedFieldNames.has(p.name))

    return (
      <div className="flex flex-col gap-4">
        {groups.filter((g: GroupInfo) =>
          g.fields.some((f: string) => properties.find((p: any) => p.name === f))
        ).map((group, idx) => (
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

  return <DefaultObjectFieldTemplate {...props} />
}

export interface ShadowFormHostProps {
  isDark: boolean
  schema: any
  uiSchema: any
  formData: Record<string, unknown>
  onChange: (formData: Record<string, unknown>) => void
  formRef?: React.RefObject<any>
  widgets?: Record<string, any>
  templates?: Record<string, any>
  onSubmit?: (formData: Record<string, unknown>) => void
}

function ShadowFormContent({
  schema,
  uiSchema,
  formData,
  onChange,
  formRef,
  widgets,
  templates,
  onSubmit,
}: Omit<ShadowFormHostProps, 'isDark'>) {
  return (
    <ShadcnForm
      ref={formRef}
      schema={schema}
      uiSchema={uiSchema ?? {}}
      formData={formData}
      onChange={(e: any) => onChange(e.formData ?? {})}
      onSubmit={(e: any) => onSubmit?.(e.formData ?? {})}
      validator={validator}
      widgets={widgets}
      templates={templates}
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

export function ShadowFormHost(props: ShadowFormHostProps) {
  const hostRef = useRef<HTMLDivElement | null>(null)
  const shadowRef = useRef<ShadowRoot | null>(null)
  const rootRef = useRef<Root | null>(null)
  const innerRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    if (!hostRef.current || shadowRef.current) return
    const existing = hostRef.current.shadowRoot
    const shadow = existing ?? hostRef.current.attachShadow({ mode: 'open' })
    shadowRef.current = shadow

    const themeStyleEl = document.createElement('style')
    themeStyleEl.setAttribute('data-rjsf-shadcn-theme', 'default')
    themeStyleEl.textContent = getShadcnDefaultSubThemeCss()
    shadow.appendChild(themeStyleEl)

    const overrideStyleEl = document.createElement('style')
    overrideStyleEl.setAttribute('data-rjsf-app-override', 'true')
    overrideStyleEl.textContent = appColorOverrideCss
    shadow.appendChild(overrideStyleEl)

    const inner = document.createElement('div')
    inner.setAttribute('data-rjsf-form-root', 'true')
    inner.style.height = '100%'
    inner.style.width = '100%'
    inner.style.color = 'var(--foreground)'
    inner.style.backgroundColor = 'var(--background)'
    inner.style.fontFamily = 'inherit'
    shadow.appendChild(inner)
    innerRef.current = inner

    rootRef.current = createRoot(inner)
    rootRef.current.render(<ShadowFormContent {...props} />)

    return () => {
      rootRef.current?.unmount()
      if (shadowRef.current) {
        shadowRef.current
          .querySelectorAll('style[data-rjsf-shadcn-theme], style[data-rjsf-app-override], div[data-rjsf-form-root]')
          .forEach((el) => el.remove())
      }
      rootRef.current = null
      shadowRef.current = null
      innerRef.current = null
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    if (hostRef.current) {
      hostRef.current.classList.toggle('dark', props.isDark)
    }
  }, [props.isDark])

  // Re-render on all prop changes
  useEffect(() => {
    rootRef.current?.render(<ShadowFormContent {...props} />)
  }, [props.schema, props.uiSchema, props.formData, props.onChange, props.onSubmit, props.widgets, props.templates])

  return <div ref={hostRef} className="h-full w-full" />
}
