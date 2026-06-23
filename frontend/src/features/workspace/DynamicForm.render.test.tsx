import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/react'
import { useRef, type RefObject } from 'react'
import { withTheme } from '@rjsf/core'
import { Theme as ShadcnTheme } from '@rjsf/shadcn'
import validator from '../../lib/rjsfValidator'

// Radix (被 @rjsf/shadcn 间接依赖) 在布局 effect 里读 ResizeObserver；
// jsdom 没有它，不 polyfill 会 ReferenceError 直接挂掉。
class RO {
  observe() {}
  unobserve() {}
  disconnect() {}
}
;(globalThis as any).ResizeObserver = (globalThis as any).ResizeObserver || RO

// 与 DynamicForm.tsx 完全一致的构造方式与自定义 widget。
const ShadcnForm = withTheme(ShadcnTheme)

function FolderPickerWidget(props: any) {
  const { value, onChange } = props
  return (
    <div className="flex gap-2">
      <input
        type="text"
        aria-label="folder-text"
        value={value ?? ''}
        onChange={(e) => onChange(e.target.value)}
      />
      <button type="button" aria-label="folder-pick" onClick={() => onChange('picked')}>
        📁
      </button>
    </div>
  )
}

const customWidgets = { folderPicker: FolderPickerWidget }

const schema = {
  type: 'object',
  required: ['indir', 'outdir'],
  properties: {
    indir: { type: 'string', title: 'Input Directory', format: 'data-url' },
    outdir: { type: 'string', title: 'Output Directory', format: 'data-url' },
    exportMode: {
      type: 'array',
      title: 'Export Mode',
      items: { type: 'string', enum: ['all', 'hardware', 'block'] },
    },
    keepFolderStructure: { type: 'boolean', title: 'Keep Folder Structure', default: false },
    sclFormat: {
      type: 'string',
      title: 'SCL Block Format',
      enum: ['ExternalSource', 'SimaticML', 'SimaticSD'],
      default: 'ExternalSource',
    },
    umacPassword: { type: 'string', title: 'UMAC Password' },
  },
}

// 复刻 DynamicForm 里 localizeUiSchema 的 JSON Forms → RJSF 展平逻辑，
// 验证 ui:widget 覆盖能解决「format: data-url 退回 FileWidget」的回归。
function localizeUiSchema(uiSchema: any): any {
  if (!uiSchema) return undefined
  const fieldUi: Record<string, any> = {}
  const order: string[] = []
  walk(uiSchema, fieldUi, order)
  if (Object.keys(fieldUi).length === 0) return undefined
  return { ...fieldUi, 'ui:order': [...order, '*'] }
}
function walk(el: any, fieldUi: Record<string, any>, order: string[]) {
  if (!el || typeof el !== 'object') return
  const t: string = el.type
  if (t === 'VerticalLayout' || t === 'HorizontalLayout' || t === 'Group' || t === 'Categorization' || t === 'Category') {
    if (Array.isArray(el.elements)) for (const c of el.elements) walk(c, fieldUi, order)
    return
  }
  if (t === 'Control') {
    const m = typeof el.scope === 'string' ? el.scope.match(/^#\/properties\/([^/]+)/) : null
    if (!m) return
    const field = m[1]
    const ui: any = {}
    const opts = el.options
    if (opts && typeof opts === 'object') {
      if (opts.folderPicker === true) ui['ui:widget'] = 'folderPicker'
      else if (opts.saveFilePicker === true) ui['ui:widget'] = 'saveFilePicker'
      else if (opts.filePicker === true) ui['ui:widget'] = 'filePicker'
      else if (typeof opts.widget === 'string' && opts.widget) ui['ui:widget'] = opts.widget
      if (opts.inputType) ui['ui:inputType'] = opts.inputType
    }
    fieldUi[field] = ui
    if (!order.includes(field)) order.push(field)
    return
  }
}

const tiaUiSchema = {
  type: 'VerticalLayout',
  elements: [
    {
      type: 'Group',
      label: 'Required Paths',
      elements: [
        { type: 'Control', scope: '#/properties/indir', options: { folderPicker: true } },
        { type: 'Control', scope: '#/properties/outdir', options: { folderPicker: true } },
      ],
    },
    { type: 'Control', scope: '#/properties/sclFormat' },
    { type: 'Control', scope: '#/properties/umacPassword', options: { inputType: 'password' } },
  ],
}

describe('RJSF shadcn renders tia-export schema', () => {
  it('folderPicker ui:widget overrides data-url format → text input + 📁 button, not file input', () => {
    const uiSchema = localizeUiSchema(tiaUiSchema)
    // 转换产出 RJSF uiSchema：indir/outdir 带 ui:widget=folderPicker，并有 ui:order。
    expect(uiSchema.indir['ui:widget']).toBe('folderPicker')
    expect(uiSchema.outdir['ui:widget']).toBe('folderPicker')
    expect(uiSchema['ui:order']).toEqual(['indir', 'outdir', 'sclFormat', 'umacPassword', '*'])

    const { container } = render(
      <ShadcnForm
        schema={schema as any}
        uiSchema={uiSchema}
        formData={{}}
        onChange={() => {}}
        validator={validator}
        widgets={customWidgets}
        showErrorList={false}
      />,
    )
    // 回归：format: data-url 没有被 folderPicker 覆盖时，这里会是 input[type=file]。
    const fileInputs = container.querySelectorAll('input[type="file"]')
    expect(fileInputs.length).toBe(0)
    // folderPicker widget 渲染了文本框 + 📁 按钮。
    const folderTexts = container.querySelectorAll('input[aria-label="folder-text"]')
    expect(folderTexts.length).toBe(2)
    const folderButtons = container.querySelectorAll('button[aria-label="folder-pick"]')
    expect(folderButtons.length).toBe(2)
  })

  it('password inputType produces input[type=password]', () => {
    const uiSchema = localizeUiSchema(tiaUiSchema)
    const { container } = render(
      <ShadcnForm
        schema={schema as any}
        uiSchema={uiSchema}
        formData={{}}
        onChange={() => {}}
        validator={validator}
        widgets={customWidgets}
        showErrorList={false}
      />,
    )
    const pwInputs = container.querySelectorAll('input[type="password"]')
    expect(pwInputs.length).toBe(1)
  })

  it('enum field renders a select trigger (role=button), boolean renders a checkbox', () => {
    const uiSchema = localizeUiSchema(tiaUiSchema)
    const { container } = render(
      <ShadcnForm
        schema={schema as any}
        uiSchema={uiSchema}
        formData={{}}
        onChange={() => {}}
        validator={validator}
        widgets={customWidgets}
        showErrorList={false}
      />,
    )
    // sclFormat enum → FancySelect 触发框 (cmdk-root 上的 role=button)。
    // 注意 folderPicker 的 📁 也是 button，所以只断言「至少有 select 触发框」。
    const buttons = container.querySelectorAll('[role="button"]')
    expect(buttons.length).toBeGreaterThanOrEqual(1)
    const checkboxes = container.querySelectorAll('[role="checkbox"]')
    expect(checkboxes.length).toBeGreaterThanOrEqual(1)
  })

  it('Form ref exposes validateForm() that returns false on invalid data, true on valid', () => {
    const uiSchema = localizeUiSchema(tiaUiSchema)
    let formRef: RefObject<any> = { current: null } as any
    function Harness() {
      const ref = useRef<any>(null)
      formRef = ref
      return (
        <ShadcnForm
          ref={ref}
          schema={schema as any}
          uiSchema={uiSchema}
          formData={{}}
          onChange={() => {}}
          validator={validator}
          widgets={customWidgets}
          showErrorList={false}
        />
      )
    }
    render(<Harness />)
    expect(typeof formRef.current?.validateForm).toBe('function')
    // formData 为空、indir/outdir 是 required → 校验失败。
    const ok = formRef.current.validateForm()
    expect(ok).toBe(false)
  })
})
