import { create } from 'zustand'
import type { OutputLine, OutputEventDto } from '../types'

// ---------- rAF batching ----------
// 收集 1 帧内的输出行，批量 commit 到 Zustand state
let pendingBatch: Record<string, OutputLine[]> = {}
let rafId: number | null = null

function flushBatch(storeSet: any) {
  rafId = null
  const batch = pendingBatch
  pendingBatch = {}
  if (Object.keys(batch).length === 0) return

  storeSet((state: OutputState) => {
    const nextOutputs = { ...state.outputs }
    for (const [tid, incoming] of Object.entries(batch)) {
      const current = nextOutputs[tid] ?? []
      const merged = [...current, ...incoming]
      if (merged.length > MAX_OUTPUT_LINES) {
        merged.splice(0, merged.length - MAX_OUTPUT_LINES)
      }
      nextOutputs[tid] = merged
    }
    return { outputs: nextOutputs }
  })
}

function scheduleFlush(set: any) {
  if (rafId === null) {
    rafId = requestAnimationFrame(() => flushBatch(set))
  }
}
// ----------------------------------

interface OutputState {
  outputs: Record<string, OutputLine[]>

  appendOutput: (taskId: string, event: OutputEventDto) => void
  appendOutputLine: (taskId: string, line: OutputLine) => void
  setOutput: (taskId: string, lines: OutputLine[]) => void
  clearOutput: (taskId: string) => void
}

const MAX_OUTPUT_LINES = 10000

export const useOutputStore = create<OutputState>((set, get) => ({
  outputs: {},

  appendOutput: (taskId: string, event: OutputEventDto) => {
    const line: OutputLine = {
      timestamp: event.timestamp,
      source: event.source,
      level: event.level,
      message: event.message,
    }
    // 放入待刷新队列
    const lines = pendingBatch[taskId] ?? []
    lines.push(line)
    pendingBatch[taskId] = lines
    scheduleFlush(set)
  },

  appendOutputLine: (taskId: string, line: OutputLine) => {
    // 也通过 batch 路径，保持一致性
    const lines = pendingBatch[taskId] ?? []
    lines.push(line)
    pendingBatch[taskId] = lines
    scheduleFlush(set)
  },

  setOutput: (taskId: string, lines: OutputLine[]) => {
    set((state) => ({
      outputs: {
        ...state.outputs,
        [taskId]: lines,
      },
    }))
  },

  clearOutput: (taskId: string) => {
    set((state) => {
      const outputs = { ...state.outputs }
      delete outputs[taskId]
      return { outputs }
    })
  },
}))
