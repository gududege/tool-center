import * as React from 'react'
import { Group, Panel, Separator } from 'react-resizable-panels'
import type { PanelImperativeHandle } from 'react-resizable-panels'
import { ThemeToggle } from '../components/ui/ThemeToggle'
import { LanguageSwitcher } from '../components/ui/LanguageSwitcher'
import { useT } from '../contexts/TranslationContext'
import { cn } from '../components/ui/cn'

interface MainLayoutProps {
  sidebar: React.ReactNode
  workspace: React.ReactNode
  sidebarCollapsed: boolean
  sidebarSize: number
  bottomPanelSize: number
  onToggleSidebar: () => void
  bottomTab: 'tasks' | 'output'
  onBottomTabChange: (tab: 'tasks' | 'output') => void
  bottomContent: React.ReactNode
  onSidebarResize?: (size: { asPercentage: number; inPixels: number }) => void
  onBottomResize?: (size: { asPercentage: number; inPixels: number }) => void
}

export function MainLayout({
  sidebar,
  workspace,
  sidebarCollapsed,
  sidebarSize,
  bottomPanelSize,
  onToggleSidebar,
  bottomTab,
  onBottomTabChange,
  bottomContent,
  onSidebarResize,
  onBottomResize,
}: MainLayoutProps) {
  const { t } = useT()
  const sidebarRef = React.useRef<PanelImperativeHandle>(null)

  React.useEffect(() => {
    const api = sidebarRef.current
    if (!api) return
    if (sidebarCollapsed && !api.isCollapsed()) {
      api.collapse()
    } else if (!sidebarCollapsed && api.isCollapsed()) {
      api.expand()
    }
  }, [sidebarCollapsed])

  return (
    <div className="flex h-screen w-screen flex-col bg-background text-foreground">
      {/* ---- Header ---- */}
      <header className="flex h-10 shrink-0 items-center justify-between border-b border-border px-4">
        <div className="flex items-center gap-3">
          <button
            onClick={onToggleSidebar}
            className={cn(
              'flex items-center gap-1.5 rounded-md px-1.5 py-1 text-xs',
              'text-neutral-500 hover:bg-neutral-100 dark:text-neutral-400 dark:hover:bg-neutral-800',
              'transition-colors',
            )}
            title={sidebarCollapsed ? t('sidebar.title') : t('sidebar.title')}
            aria-label={sidebarCollapsed ? t('sidebar.title') : t('sidebar.title')}
          >
            <svg
              className="h-4 w-4 transition-transform"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              style={{ transform: sidebarCollapsed ? 'scaleX(-1)' : 'none' }}
            >
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
              <line x1="9" y1="3" x2="9" y2="21" />
            </svg>
            <span className="hidden sm:inline text-xs">{t('sidebar.title')}</span>
          </button>
          <h1 className="text-sm font-semibold">{t('app.title')}</h1>
        </div>
        <div className="flex items-center gap-2">
          <LanguageSwitcher />
          <ThemeToggle />
        </div>
      </header>

      {/* ---- Main ---- */}
      <Group className="flex-1 overflow-hidden">
        {/* Sidebar — collapsible */}
        <Panel
          panelRef={sidebarRef}
          defaultSize={sidebarSize}
          minSize={5}
          collapsible
          collapsedSize={0}
          onResize={onSidebarResize}
        >
          <div className="h-full overflow-hidden border-r-[3px] border-neutral-300 dark:border-neutral-700">
            {sidebar}
          </div>
        </Panel>

        {/* Sidebar resize handle — hidden when collapsed */}
        {!sidebarCollapsed && (
          <Separator
            className={cn(
              'group relative w-1.5 shrink-0 cursor-col-resize',
              'flex items-center justify-center',
              'bg-neutral-300/60 dark:bg-neutral-700/60',
              'transition-colors',
              'hover:bg-blue-400/30 data-[resize-handle-active]:bg-blue-400/50',
            )}
          >
            <div className="h-6 w-0.5 rounded-full bg-neutral-400/50 transition-opacity group-hover:bg-blue-400 dark:bg-neutral-500/50" />
          </Separator>
        )}

        {/* Workspace + Bottom panel */}
        <Panel defaultSize={82} minSize={20}>
          <Group orientation="vertical">
            {/* Parameters — always visible */}
            <Panel defaultSize={65} minSize={10}>
              <div className="h-full overflow-hidden">{workspace}</div>
            </Panel>

            {/* Bottom resize handle with thick line */}
            <Separator
              className={cn(
                'group relative h-1.5 shrink-0 cursor-row-resize',
                'flex items-center justify-center',
                'bg-neutral-300/60 dark:bg-neutral-700/60',
                'transition-colors',
                'hover:bg-blue-400/30 data-[resize-handle-active]:bg-blue-400/50',
              )}
            >
              <div className="w-6 h-0.5 rounded-full bg-neutral-400/50 transition-opacity group-hover:bg-blue-400 dark:bg-neutral-500/50" />
            </Separator>

            {/* Bottom panel: Tasks | Output (tabbed) */}
            <Panel
              defaultSize={bottomPanelSize}
              minSize={5}
              onResize={onBottomResize}
            >
              <div className="flex h-full flex-col overflow-hidden">
                {/* Tab bar */}
                <div className="flex shrink-0 items-center border-b border-border bg-neutral-50 dark:bg-neutral-900">
                  <button
                    className={cn(
                      'flex items-center gap-1.5 px-4 py-1.5 text-xs font-medium transition-colors',
                      'border-r border-border',
                      bottomTab === 'tasks'
                        ? 'bg-background text-foreground'
                        : 'text-muted-foreground hover:bg-neutral-100 dark:hover:bg-neutral-800',
                    )}
                    onClick={() => onBottomTabChange('tasks')}
                  >
                    <svg className="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
                      <line x1="9" y1="3" x2="9" y2="21" />
                      <line x1="3" y1="9" x2="21" y2="9" />
                    </svg>
                    {t('panel.tasks')}
                  </button>
                  <button
                    className={cn(
                      'flex items-center gap-1.5 px-4 py-1.5 text-xs font-medium transition-colors',
                      'border-r border-border',
                      bottomTab === 'output'
                        ? 'bg-background text-foreground'
                        : 'text-muted-foreground hover:bg-neutral-100 dark:hover:bg-neutral-800',
                    )}
                    onClick={() => onBottomTabChange('output')}
                  >
                    <svg className="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <polyline points="16 18 22 12 16 6" />
                      <polyline points="8 6 2 12 8 18" />
                    </svg>
                    {t('panel.output')}
                  </button>
                </div>
                {/* Tab content */}
                <div className="flex-1 overflow-hidden">
                  {bottomContent}
                </div>
              </div>
            </Panel>
          </Group>
        </Panel>
      </Group>
    </div>
  )
}