import { create } from 'zustand'
import type { OutputLine, OutputEventDto } from '../types'

interface OutputState {
  outputs: Record<string, OutputLine[]>

  appendOutput: (taskId: string, event: OutputEventDto) => void
  appendOutputLine: (taskId: string, line: OutputLine) => void
  setOutput: (taskId: string, lines: OutputLine[]) => void
  clearOutput: (taskId: string) => void
  getOutput: (taskId: string) => OutputLine[]
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
    get().appendOutputLine(taskId, line)
  },

  appendOutputLine: (taskId: string, line: OutputLine) => {
    set((state) => {
      const current = state.outputs[taskId] ?? []
      const next = [...current, line]
      // FIFO 裁剪
      if (next.length > MAX_OUTPUT_LINES) {
        next.splice(0, next.length - MAX_OUTPUT_LINES)
      }
      return {
        outputs: {
          ...state.outputs,
          [taskId]: next,
        },
      }
    })
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

  getOutput: (taskId: string) => {
    return get().outputs[taskId] ?? []
  },
}))
