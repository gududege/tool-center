import { useState, useEffect } from 'react'
import { MainLayout } from './layouts/MainLayout'
import { Sidebar } from './features/sidebar'
import { Workspace } from './features/workspace'
import { TaskPanel } from './features/taskPanel'
import { OutputViewer } from './features/outputPanel'
import { useTaskEvents } from './hooks/useTaskEvents'
import { useSettingsStore } from './stores/settingsStore'
import { useT } from './contexts/TranslationContext'

export default function App() {
  const settings = useSettingsStore((s) => s.settings)
  const fetchSettings = useSettingsStore((s) => s.fetchSettings)
  const updateLayout = useSettingsStore((s) => s.updateLayout)
  const setLanguage = useSettingsStore((s) => s.setLanguage)
  const { t, lang, setLanguage: setUILanguage } = useT()

  const [sidebarCollapsed, setSidebarCollapsed] = useState(false)
  const [bottomTab, setBottomTab] = useState<'tasks' | 'output'>('tasks')
  const [outputTaskId, setOutputTaskId] = useState<string | null>(null)

  // 初始化: 加载设置
  useEffect(() => {
    fetchSettings()
  }, [])

  // 设置加载后初始化布局状态
  useEffect(() => {
    if (!settings) return
    setSidebarCollapsed(settings.sidebarCollapsed)
    setBottomTab(settings.bottomTab as 'tasks' | 'output')
    if (settings.language && settings.language !== lang) {
      setUILanguage(settings.language)
    }
  }, [settings])

  // 语言切换时同步到后端设置
  useEffect(() => {
    if (!settings || lang === settings.language) return
    setLanguage(lang)
  }, [lang])

  // 注册全局事件监听
  useTaskEvents()

  // 侧边栏切换 → 保存
  const handleToggleSidebar = () => {
    const next = !sidebarCollapsed
    setSidebarCollapsed(next)
    updateLayout({ sidebarCollapsed: next })
  }

  // 标签切换 → 保存
  const handleBottomTabChange = (tab: 'tasks' | 'output') => {
    setBottomTab(tab)
    updateLayout({ bottomTab: tab })
  }

  // 面板大小变化 → 保存
  const handleSidebarResize = (size: { asPercentage: number }) => {
    if (size.asPercentage > 0) {
      updateLayout({ sidebarSize: Math.round(size.asPercentage) })
    }
  }
  const handleBottomResize = (size: { asPercentage: number }) => {
    updateLayout({ bottomPanelSize: Math.round(size.asPercentage) })
  }

  // 点击任务时切换到 output tab 并显示其输出
  const handleTaskClick = (taskId: string) => {
    setOutputTaskId(taskId)
    setBottomTab('output')
  }

  // 底部面板内容
  const bottomContent =
    bottomTab === 'tasks' ? (
      <TaskPanel onTaskClick={handleTaskClick} />
    ) : outputTaskId ? (
      <OutputViewer key={outputTaskId} taskId={outputTaskId} />
    ) : (
      <div className="flex h-full items-center justify-center text-xs text-neutral-400">
        {t('panel.select_task')}
      </div>
    )

  const sidebarSize = settings?.sidebarSize ?? 18
  const bottomPanelSize = settings?.bottomPanelSize ?? 35

  // 等待设置加载后再渲染布局，确保面板尺寸正确恢复
  if (!settings) {
    return (
      <div className="flex h-screen w-screen items-center justify-center bg-background text-sm text-muted-foreground">
        {t('app.loading')}
      </div>
    )
  }

  return (
    <MainLayout
      sidebar={<Sidebar />}
      workspace={<Workspace />}
      sidebarCollapsed={sidebarCollapsed}
      sidebarSize={sidebarSize}
      bottomPanelSize={bottomPanelSize}
      onToggleSidebar={handleToggleSidebar}
      bottomTab={bottomTab}
      onBottomTabChange={handleBottomTabChange}
      bottomContent={bottomContent}
      onSidebarResize={handleSidebarResize}
      onBottomResize={handleBottomResize}
    />
  )
}
