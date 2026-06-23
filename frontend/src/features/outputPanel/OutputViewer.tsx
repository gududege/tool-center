import { useRef, useEffect, useState } from 'react'
import { useOutputStore } from '../../stores/outputStore'
import { ScrollArea } from '../../components/ui/ScrollArea'
import { Button } from '../../components/ui/Button'
import { useT } from '../../contexts/TranslationContext'
import type { OutputLine } from '../../types'

interface OutputViewerProps {
  taskId: string
}

export function OutputViewer({ taskId }: OutputViewerProps) {
  const { t } = useT()
  const outputs = useOutputStore((s) => s.outputs)
  const lines = outputs[taskId] ?? []
  const [autoScroll, setAutoScroll] = useState(true)
  const bottomRef = useRef<HTMLDivElement>(null)
  const [search, setSearch] = useState('')

  // 自动滚动
  useEffect(() => {
    if (autoScroll && bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [lines.length, autoScroll])

  const filteredLines: OutputLine[] = search
    ? lines.filter((l) => l.message.toLowerCase().includes(search.toLowerCase()))
    : lines

  const handleCopy = () => {
    const text = lines.map((l) => `[${l.timestamp}] [${l.source}] ${l.message}`).join('\n')
    navigator.clipboard.writeText(text)
  }

  return (
    <div className="flex h-full flex-col">
      {/* 工具栏 */}
      <div className="flex items-center gap-2 border-b border-neutral-200 px-3 py-1.5 dark:border-neutral-800">
        <span className="text-xs font-semibold text-neutral-500">{t('panel.output')}</span>
        <input
          type="text"
          className="ml-2 h-6 flex-1 rounded border border-neutral-300 bg-transparent px-2 text-xs dark:border-neutral-700"
          placeholder={t('common.search')}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setAutoScroll(!autoScroll)}
        >
          {autoScroll ? t('output.auto_scroll') : t('output.manual_scroll')}
        </Button>
        <Button variant="ghost" size="sm" onClick={handleCopy}>
          {t('output.copy')}
        </Button>
      </div>

      {/* 日志 */}
      <ScrollArea className="flex-1">
        <div className="p-2 font-mono text-xs leading-5">
          {filteredLines.length === 0 ? (
            <p className="text-neutral-400">{t('output.no_output')}</p>
          ) : (
            filteredLines.map((line, i) => (
              <div
                key={i}
                className={`whitespace-pre-wrap ${
                  line.source === 'stderr'
                    ? 'text-red-400'
                    : line.source === 'system'
                      ? 'text-blue-400'
                      : 'text-neutral-300'
                }`}
              >
                <span className="text-neutral-600">
                  [{line.timestamp?.slice(11, 19) ?? '--:--:--'}]
                </span>{' '}
                {line.message}
              </div>
            ))
          )}
          <div ref={bottomRef} />
        </div>
      </ScrollArea>
    </div>
  )
}
