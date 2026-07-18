import anser from 'anser'

const ansiCache = new Map<string, string>()

export function ansiToHtml(text: string): string {
  // 快速路径：没有 ANSI 码直接返回
  if (!text.includes('\x1b[')) return text

  const cached = ansiCache.get(text)
  if (cached !== undefined) return cached

  const html = anser.ansiToHtml(anser.escapeForHtml(text), { use_classes: true })
  // 简单 LRU：超 500 条时清理
  if (ansiCache.size > 500) ansiCache.clear()
  ansiCache.set(text, html)
  return html
}
