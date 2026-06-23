import { describe, it, expect } from 'vitest'
import { buildMenuTree } from './Sidebar'
import { createMockPlugin, createMockPlugins } from '../../test/mocks'
import type { MenuNode } from '../../types'

describe('buildMenuTree', () => {
  it('returns empty array for empty plugin list', () => {
    expect(buildMenuTree([])).toEqual([])
  })

  it('creates a flat tree for a plugin with no navigation groups', () => {
    const plugin = createMockPlugin({ navigation: { group: [], order: 0 } })
    const tree = buildMenuTree([plugin])
    expect(tree).toHaveLength(1)
    expect(tree[0]).toMatchObject({
      id: 'plugin-tia-export',
      name: 'TIA Portal Export',
      pluginId: 'tia-export',
    })
  })

  it('builds nested groups from navigation.group array', () => {
    const plugin = createMockPlugin({
      navigation: { group: ['Industrial Automation', 'TIA Portal'], order: 1 },
    })
    const tree = buildMenuTree([plugin])

    // root → [Industrial Automation]
    expect(tree).toHaveLength(1)
    const group1 = tree[0]
    expect(group1).toMatchObject({
      id: 'group-Industrial Automation',
      name: 'Industrial Automation',
    })
    expect(group1.pluginId).toBeUndefined()

    // Industrial Automation → [TIA Portal]
    expect(group1.children).toHaveLength(1)
    const group2 = group1.children[0]
    expect(group2).toMatchObject({
      id: 'group-Industrial Automation/TIA Portal',
      name: 'TIA Portal',
    })

    // TIA Portal → [plugin leaf]
    expect(group2.children).toHaveLength(1)
    expect(group2.children[0]).toMatchObject({
      id: 'plugin-tia-export',
      name: 'TIA Portal Export',
      pluginId: 'tia-export',
    })
  })

  it('groups plugins under same navigation group', () => {
    const plugins = [
      createMockPlugin({ id: 'p1', name: 'Plugin 1', navigation: { group: ['Database'], order: 1 } }),
      createMockPlugin({ id: 'p2', name: 'Plugin 2', navigation: { group: ['Database'], order: 2 } }),
    ]
    const tree = buildMenuTree(plugins)

    expect(tree).toHaveLength(1)
    expect(tree[0].name).toBe('Database')
    expect(tree[0].children).toHaveLength(2)
    expect(tree[0].children[0].pluginId).toBe('p1')
    expect(tree[0].children[1].pluginId).toBe('p2')
  })

  it('handles multiple top-level groups', () => {
    const plugins = createMockPlugins(4)
    const tree = buildMenuTree(plugins)

    // Even-indexed plugins go to "Group A", odd to "Group B > Subgroup"
    const groupA = tree.find((n) => n.name === 'Group A')
    const groupB = tree.find((n) => n.name === 'Group B')
    expect(groupA).toBeDefined()
    expect(groupB).toBeDefined()
    expect(groupA!.children).toHaveLength(2)
    expect(groupB!.children).toHaveLength(1)
    expect(groupB!.children[0].name).toBe('Subgroup')
    expect(groupB!.children[0].children).toHaveLength(2)
  })

  it('preserves plugin order within a group', () => {
    const plugins = [
      createMockPlugin({ id: 'z', name: 'Z Last', navigation: { group: ['Tools'], order: 10 } }),
      createMockPlugin({ id: 'a', name: 'A First', navigation: { group: ['Tools'], order: 1 } }),
    ]
    const tree = buildMenuTree(plugins)
    // Note: buildMenuTree doesn't sort — it preserves insertion order
    expect(tree[0].children[0].pluginId).toBe('z')
    expect(tree[0].children[1].pluginId).toBe('a')
  })

  it('handles real Wails backend data shape', () => {
    // Exact shape returned by the Go backend via JSON serialization
    const rawData = [
      {
        id: 'tia-export',
        name: 'TIA Portal Export',
        description: 'Export content from Siemens TIA Portal projects',
        version: '1.0.0',
        icon: '',
        navigation: { group: ['Industrial Automation', 'TIA Portal'], order: 1 },
      },
    ]
    const tree = buildMenuTree(rawData as any)

    expect(tree).toHaveLength(1)
    expect(tree[0].name).toBe('Industrial Automation')
    expect(tree[0].children[0].children[0].pluginId).toBe('tia-export')
  })

  it('uses Chinese groups and names when lang=zh', () => {
    const plugin = createMockPlugin({
      name: 'My Tool',
      name_cn: '我的工具',
      navigation: {
        group: ['Network', 'Ping'],
        group_cn: ['网络', 'Ping'],
        order: 1,
      },
    })
    const tree = buildMenuTree([plugin], 'zh')

    expect(tree).toHaveLength(1)
    expect(tree[0].name).toBe('网络')
    expect(tree[0].children[0].name).toBe('Ping')
    expect(tree[0].children[0].children[0].name).toBe('我的工具')
  })

  it('falls back to English when lang=zh but no _cn fields exist', () => {
    // PluginSummaryDto with no _cn fields
    const plugin = createMockPlugin({
      name_cn: undefined,
      navigation: { group: ['Network', 'Ping'], group_cn: undefined, order: 1 } as any,
    })
    const tree = buildMenuTree([plugin], 'zh')

    expect(tree[0].name).toBe('Network')
    expect(tree[0].children[0].children[0].name).toBe('TIA Portal Export')
  })
})
