import type { PluginDto } from '../types'

export function buildCommand(plugin: PluginDto | null, formData: Record<string, unknown>): string {
  if (!plugin) return ''
  const exe = plugin.execution?.target ?? ''
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
