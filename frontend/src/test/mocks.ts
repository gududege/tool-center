import type { PluginSummaryDto, NavigationDto } from '../types'

export function createMockPlugin(
  overrides: Partial<PluginSummaryDto> = {},
): PluginSummaryDto {
  return {
    id: 'tia-export',
    name: 'TIA Portal Export',
    name_cn: 'TIA Portal 导出',
    description: 'Export content from Siemens TIA Portal projects',
    description_cn: '从西门子 TIA Portal 项目导出内容',
    version: '1.0.0',
    icon: '',
    navigation: {
      group: ['Industrial Automation', 'TIA Portal'],
      group_cn: ['工业自动化', 'TIA Portal'],
      order: 1,
    },
    ...overrides,
  }
}

export function createMockPlugins(count: number = 3): PluginSummaryDto[] {
  return Array.from({ length: count }, (_, i) =>
    createMockPlugin({
      id: `plugin-${i}`,
      name: `Plugin ${i}`,
      navigation: {
        group: i % 2 === 0 ? ['Group A'] : ['Group B', 'Subgroup'],
        group_cn: i % 2 === 0 ? ['分组 A'] : ['分组 B', '子分组'],
        order: i,
      },
    }),
  )
}

/** Represents the exact JSON shape that Wails Go backend sends */
export const RAW_WAILS_PLUGIN_RESPONSE = [
  {
    id: 'tia-export',
    name: 'TIA Portal Export',
    name_cn: 'TIA Portal 导出',
    description: 'Export content from Siemens TIA Portal projects',
    description_cn: '从西门子 TIA Portal 项目导出内容',
    version: '1.0.0',
    icon: '',
    navigation: {
      group: ['Industrial Automation', 'TIA Portal'],
      group_cn: ['工业自动化', 'TIA Portal'],
      order: 1,
    },
  },
]
