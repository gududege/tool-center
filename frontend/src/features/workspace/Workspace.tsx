import { useTabStore } from '../../stores/tabStore'
import { useOutputStore } from '../../stores/outputStore'
import { DynamicForm } from './DynamicForm'
import { ScrollArea } from '../../components/ui/ScrollArea'
import { useT } from '../../contexts/TranslationContext'
import type { Tab } from '../../types'

export function Workspace() {
  const { t } = useT()
  const { tabs, activeTabId, activateTab, closeTab } = useTabStore()

  if (tabs.length === 0) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <div className="text-3xl mb-2">🔧</div>
          <p className="text-sm text-neutral-400">
            {t('workspace.empty_hint')}
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex h-full flex-col">
      {/* Tab bar */}
      <div className="flex items-center border-b border-neutral-200 bg-neutral-50 dark:border-neutral-800 dark:bg-neutral-900">
          <div className="flex overflow-x-auto">
            {tabs.map((tab) => (
              <div
                key={tab.id}
                className={`
                  group flex items-center gap-1 px-3 py-2 text-xs cursor-pointer select-none
                  border-r border-neutral-200 dark:border-neutral-800
                  transition-colors
                  ${activeTabId === tab.id
                    ? 'bg-white text-neutral-900 dark:bg-neutral-950 dark:text-neutral-100'
                    : 'text-neutral-500 hover:bg-neutral-100 dark:hover:bg-neutral-800'
                  }
                `}
                onClick={() => activateTab(tab.id)}
              >
                <span>{tab.title}</span>
                {tab.dirty && <span className="h-1.5 w-1.5 rounded-full bg-blue-500" />}
                <button
                  className="ml-1 opacity-0 group-hover:opacity-100 hover:text-red-500"
                  onClick={(e) => {
                    e.stopPropagation()
                    closeTab(tab.id)
                  }}
                >
                  ×
                </button>
              </div>
            ))}
          </div>
        </div>

      {/* Tab content */}
      <div className="flex-1 overflow-hidden">
        {tabs.map((tab) => (
          <div
            key={tab.id}
            className={`h-full ${activeTabId === tab.id ? 'block' : 'hidden'}`}
          >
            {tab.taskId && !tab.pluginId.startsWith('output-') ? (
              <DynamicForm tab={tab} />
            ) : tab.taskId ? (
              <OutputTab tab={tab} />
            ) : (
              <DynamicForm tab={tab} />
            )}
          </div>
        ))}
      </div>
    </div>
  )
}

/**
 * 输出 Tab — 只读日志查看
 */
function OutputTab({ tab }: { tab: Tab }) {
  const { t } = useT()
  const outputs = useOutputStore((s) => s.outputs)
  const outputLines = tab.taskId ? outputs[tab.taskId] ?? [] : []

  return (
    <div className="flex h-full flex-col p-4">
      <h3 className="mb-2 text-sm font-medium">{t('panel.output')}: {tab.title}</h3>
      <ScrollArea className="flex-1 rounded-md bg-neutral-950 p-3 font-mono text-xs text-green-400">
        {outputLines.length === 0 ? (
          <p className="text-neutral-500">{t('output.no_output')}</p>
        ) : (
          outputLines.map((line: any, i: number) => (
            <div key={i} className="whitespace-pre-wrap">
              <span className="text-neutral-500">[{line.timestamp?.slice(11, 19) ?? '--:--:--'}]</span>{' '}
              <span className={line.level === 'error' || line.level === 'warn' ? 'text-red-400' : ''}>
                {line.message}
              </span>
            </div>
          ))
        )}
      </ScrollArea>
    </div>
  )
}
