import { useEffect, useState } from 'react'
import { usePluginStore } from '../../stores/pluginStore'
import { useTabStore } from '../../stores/tabStore'
import { ScrollArea } from '../../components/ui/ScrollArea'
import { cn } from '../../components/ui/cn'
import { useT } from '../../contexts/TranslationContext'
import type { MenuNode, PluginSummaryDto } from '../../types'

export function buildMenuTree(plugins: PluginSummaryDto[], lang: string = 'en'): MenuNode[] {
  const root: MenuNode[] = []
  const map = new Map<string, MenuNode>()

  for (const plugin of plugins) {
    const groups = lang === 'zh' && plugin.navigation.group_cn?.length
      ? plugin.navigation.group_cn
      : plugin.navigation.group
    const pluginName = lang === 'zh' && plugin.name_cn ? plugin.name_cn : plugin.name

    if (!groups || groups.length === 0) {
      root.push({
        id: `plugin-${plugin.id}`,
        name: pluginName,
        children: [],
        pluginId: plugin.id,
      })
      continue
    }

    let currentLevel = root
    for (let i = 0; i < groups.length; i++) {
      const groupName = groups[i]
      const pathKey = groups.slice(0, i + 1).join('/')
      let node = map.get(pathKey)

      if (!node) {
        node = {
          id: `group-${pathKey}`,
          name: groupName,
          children: [],
        }
        map.set(pathKey, node)
        currentLevel.push(node)
      }
      currentLevel = node.children
    }

    // 叶子节点：插件
    currentLevel.push({
      id: `plugin-${plugin.id}`,
      name: pluginName,
      children: [],
      pluginId: plugin.id,
    })
  }

  return root
}

interface TreeNodeProps {
  node: MenuNode
  depth: number
  onSelect: (pluginId: string) => void
}

function TreeNode({ node, depth, onSelect }: TreeNodeProps) {
  const [expanded, setExpanded] = useState(true)
  const hasChildren = node.children.length > 0

  const handleClick = () => {
    if (node.pluginId) {
      onSelect(node.pluginId)
    } else {
      setExpanded(!expanded)
    }
  }

  return (
    <div>
      <button
        className={cn(
          'flex w-full items-center gap-1 px-2 py-1 text-left text-sm hover:bg-neutral-100 dark:hover:bg-neutral-800',
          'transition-colors',
        )}
        style={{ paddingLeft: `${8 + depth * 16}px` }}
        onClick={handleClick}
      >
        {hasChildren && (
          <span className="text-xs text-neutral-400">
            {expanded ? '▼' : '▶'}
          </span>
        )}
        <span className={cn(node.pluginId ? 'font-normal' : 'font-medium text-neutral-600 dark:text-neutral-400')}>
          {node.name}
        </span>
      </button>
      {hasChildren && expanded && (
        <div>
          {node.children.map((child) => (
            <TreeNode
              key={child.id}
              node={child}
              depth={depth + 1}
              onSelect={onSelect}
            />
          ))}
        </div>
      )}
    </div>
  )
}

export function Sidebar() {
  const { t, lang } = useT()
  const { plugins, fetchPlugins } = usePluginStore()
  const { addTab } = useTabStore()
  const [menuTree, setMenuTree] = useState<MenuNode[]>([])

  useEffect(() => {
    fetchPlugins()
  }, [])

  useEffect(() => {
    setMenuTree(buildMenuTree(plugins, lang))
  }, [plugins, lang])

  const handlePluginSelect = (pluginId: string) => {
    const plugin = plugins.find((p) => p.id === pluginId)
    if (!plugin) return

    const pluginName = lang === 'zh' && plugin.name_cn ? plugin.name_cn : plugin.name

    addTab({
      id: '',
      pluginId: plugin.id,
      title: pluginName,
      dirty: false,
    })
  }

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between border-b border-neutral-200 px-3 py-2 dark:border-neutral-800">
        <span className="text-xs font-semibold uppercase text-neutral-500">{t('sidebar.title')}</span>
        <button
          className="text-xs text-neutral-400 hover:text-neutral-600 dark:hover:text-neutral-300"
          onClick={() => usePluginStore.getState().reloadPlugins()}
          title={t('sidebar.reload')}
        >
          ↻
        </button>
      </div>
      <ScrollArea className="flex-1">
        {menuTree.length === 0 ? (
          <div className="p-4 text-center text-sm text-neutral-400">
            {t('sidebar.no_plugins')}
          </div>
        ) : (
          menuTree.map((node) => (
            <TreeNode
              key={node.id}
              node={node}
              depth={0}
              onSelect={handlePluginSelect}
            />
          ))
        )}
      </ScrollArea>
    </div>
  )
}
