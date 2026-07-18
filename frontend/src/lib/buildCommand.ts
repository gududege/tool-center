import type { PluginDto } from '../types'

export function buildCommand(
  plugin: PluginDto | null,
  formData: Record<string, unknown>,
  runnerPaths?: Record<string, string>,
  runnerOptions?: Record<string, string>,
): string {
  if (!plugin) return ''
  const exe = plugin.execution?.target ?? ''
  const runner = plugin.execution?.runner ?? ''
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
          for (const v of val) {
            const s = String(v)
            if (s) args.push(s)
          }
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
            const s = String(v)
            if (s) args.push(p.flag ?? `--${p.field}`, s)
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

  const parts: string[] = []
  if (runner === 'binary' || runner === '') {
    parts.push(exe)
  } else {
    const runnerPath = runnerPaths?.[runner] ?? runner
    parts.push(runnerPath)
    const opts = runnerOptions?.[runner] ?? ''
    for (const opt of opts.split(' ')) {
      if (opt) parts.push(opt)
    }
    parts.push(exe)
  }
  parts.push(...args)
  return parts.filter(Boolean).join(' ')
}
