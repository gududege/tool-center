import type { GroupInfo } from '../components/ui/ShadowFormHost'

export function localizeSchema(schema: any, lang: string): any {
  if (!schema || !schema.properties) return schema
  const out = { ...schema }
  if (schema.title_cn && lang === 'zh') out.title = schema.title_cn
  if (schema.description_cn && lang === 'zh') out.description = schema.description_cn

  const props: Record<string, unknown> = {}
  for (const [key, val] of Object.entries(schema.properties)) {
    const prop = { ...(val as any) }
    if (lang === 'zh' && prop.title_cn) prop.title = prop.title_cn
    if (lang === 'zh' && prop.description_cn) prop.description = prop.description_cn
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
    props[key] = prop
  }
  out.properties = props
  return out
}

export function localizeUiSchema(uiSchema: any, lang: string): any {
  if (!uiSchema) return undefined
  const fieldUi: Record<string, unknown> = {}
  const order: string[] = []
  const groups: GroupInfo[] = []
  walkJsonFormsLayout(uiSchema, fieldUi, order, groups, lang)
  if (Object.keys(fieldUi).length === 0 && groups.length === 0) return undefined
  const out: Record<string, unknown> = { ...fieldUi, 'ui:order': [...order, '*'] }
  if (groups.length > 0) {
    out['ui:groups'] = groups
  }
  return out
}

function walkJsonFormsLayout(
  el: any,
  fieldUi: Record<string, unknown>,
  order: string[],
  groups: GroupInfo[],
  lang: string,
) {
  if (!el || typeof el !== 'object') return
  const type: string = el.type

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

  if (type === 'VerticalLayout' || type === 'HorizontalLayout' || type === 'Categorization' || type === 'Category') {
    if (Array.isArray(el.elements)) {
      for (const child of el.elements) walkJsonFormsLayout(child, fieldUi, order, groups, lang)
    }
    return
  }

  collectControlField(el, fieldUi, order, lang)
}

function collectControlField(
  el: any,
  fieldUi: Record<string, unknown>,
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
      ui['ui:widget'] = opts.widget
    }
    if (opts.inputType) {
      ui['ui:inputType'] = opts.inputType
    }
  }
  fieldUi[field] = ui
  if (!order.includes(field)) order.push(field)
  return field
}

function scopeToField(scope: any): string | null {
  if (typeof scope !== 'string') return null
  const m = scope.match(/^#\/properties\/([^/]+)/)
  return m ? m[1] : null
}
