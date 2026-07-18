import React, { useRef, useEffect, useState, useMemo } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'
import { useOutputStore } from '../../stores/outputStore'
import { Button } from '../../components/ui/Button'
import { useT } from '../../contexts/TranslationContext'
import { ansiToHtml } from '../../lib/ansiToHtml'
import type { OutputLine } from '../../types'

interface OutputViewerProps {
  taskId: string
}

// 单行渲染组件，利用 React.memo 避免不必要的重渲染
const OutputLineRow = React.memo(function OutputLineRow({
  line,
}: {
  line: OutputLine
}) {
  return (
    <div
      className={`whitespace-pre-wrap ${
        line.source === 'stderr'
          ? 'text-red-600 dark:text-red-400'
          : line.source === 'system'
            ? 'text-blue-600 dark:text-blue-400'
            : 'text-neutral-800 dark:text-neutral-300'
      }`}
    >
      <span className="text-neutral-400 dark:text-neutral-600">
        [{line.timestamp?.slice(11, 19) ?? '--:--:--'}]
      </span>{' '}
      <span dangerouslySetInnerHTML={{ __html: ansiToHtml(line.message) }} />
    </div>
  )
})

export function OutputViewer({ taskId }: OutputViewerProps) {
  const { t } = useT()
  // 窄 selector：只订阅当前 task 的输出数组
  const lines = useOutputStore((s) => s.outputs[taskId] ?? [])
  const [autoScroll, setAutoScroll] = useState(true)
  const [search, setSearch] = useState('')
  const parentRef = useRef<HTMLDivElement>(null)
  // 记录滚动前是否在底部，用于自动锁定
  const isNearBottomRef = useRef(true)

  const filteredLines = useMemo(
    () =>
      search
        ? lines.filter((l) => l.message.toLowerCase().includes(search.toLowerCase()))
        : lines,
    [lines, search],
  )

  // 行高估算（可根据实际行高调整）
  const ROW_HEIGHT = 24

  const virtualizer = useVirtualizer({
    count: filteredLines.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => ROW_HEIGHT,
    overscan: 10,
  })

  // 滚动检测：用户手动滚上去时暂停自动滚动
  const handleScroll = () => {
    const el = parentRef.current
    if (!el) return
    const threshold = 50
    isNearBottomRef.current = el.scrollHeight - el.scrollTop - el.clientHeight < threshold
    if (isNearBottomRef.current && !autoScroll) {
      setAutoScroll(true)
    } else if (!isNearBottomRef.current && autoScroll) {
      setAutoScroll(false)
    }
  }

  // 新行到来时自动滚动到末尾
  const prevLengthRef = useRef(0)
  useEffect(() => {
    if (autoScroll && filteredLines.length > prevLengthRef.current) {
      const el = parentRef.current
      if (el) {
        el.scrollTop = el.scrollHeight
      }
    }
    prevLengthRef.current = filteredLines.length
  })

  const handleCopy = () => {
    const text = lines.map((l) => `[${l.timestamp}] [${l.source}] ${l.message}`).join('\n')
    navigator.clipboard.writeText(text).catch(() => {})
  }

  return (
    <div className="flex h-full flex-col">
      {/* 工具栏 */}
      <div className="flex items-center gap-2 border-b border-neutral-200 px-3 py-1.5 dark:border-neutral-800">
        <span className="text-xs font-semibold text-neutral-500">{t('panel.output')}</span>
        <input
          type="text"
          className="ml-2 h-6 flex-1 rounded border border-neutral-300 bg-transparent px-2 text-xs dark:border-neutral-700"
          placeholder={t('common.search') ?? 'Search'}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setAutoScroll(!autoScroll)}
          title={t('button.auto_scroll') ?? 'Auto scroll'}
        >
          {autoScroll ? t('output.auto_scroll') ?? 'Auto' : t('output.manual_scroll') ?? 'Manual'}
        </Button>
        <Button variant="ghost" size="sm" onClick={handleCopy} title={t('button.copy_output') ?? 'Copy'}>
          {t('output.copy') ?? 'Copy'}
        </Button>
      </div>

      {/* 虚拟滚动容器 */}
      <div
        ref={parentRef}
        onScroll={handleScroll}
        className="flex-1 overflow-auto"
        style={{ contain: 'strict' }}
      >
        {filteredLines.length === 0 ? (
          <div className="p-2 font-mono text-xs leading-5">
            <p className="text-neutral-400">{t('output.no_output')}</p>
          </div>
        ) : (
          <div
            className="relative p-2 font-mono text-xs leading-5"
            style={{ height: `${virtualizer.getTotalSize()}px` }}
          >
            {virtualizer.getVirtualItems().map((virtualRow) => {
              const line = filteredLines[virtualRow.index]
              return (
                <div
                  key={virtualRow.key}
                  className="absolute left-0 w-full"
                  style={{
                    height: `${virtualRow.size}px`,
                    transform: `translateY(${virtualRow.start}px)`,
                  }}
                >
                  <OutputLineRow line={line} />
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}
