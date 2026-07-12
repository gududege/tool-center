import { useEffect } from 'react'
import { useTaskStore } from '../../stores/taskStore'
import { useOutputStore } from '../../stores/outputStore'
import { taskApi } from '../../services/taskApi'
import { ScrollArea } from '../../components/ui/ScrollArea'
import { Button } from '../../components/ui/Button'
import { useT } from '../../contexts/TranslationContext'
import type { TaskDto } from '../../types'

const STATUS_COLORS: Record<string, string> = {
  created: 'bg-neutral-300 dark:bg-neutral-500',
  queued: 'bg-amber-500 dark:bg-amber-400',
  running: 'bg-blue-500 dark:bg-blue-400',
  completed: 'bg-emerald-500 dark:bg-emerald-400',
  failed: 'bg-red-500 dark:bg-red-400',
  cancelled: 'bg-orange-500 dark:bg-orange-400',
}

interface TaskItemProps {
  task: TaskDto
  onTaskClick: (taskId: string) => void
}

function TaskItem({ task, onTaskClick }: TaskItemProps) {
  const { t } = useT()

  const handleOpenOutput = () => {
    onTaskClick(task.id)
  }

  const handleCancel = async () => {
    await taskApi.cancelTask(task.id)
  }

  const handleDelete = async () => {
    await taskApi.deleteTask(task.id)
    useTaskStore.getState().removeTask(task.id)
    useOutputStore.getState().clearOutput(task.id)
  }

  const duration = task.startedAt && task.endedAt
    ? Math.round(
        (new Date(task.endedAt).getTime() - new Date(task.startedAt).getTime()) / 1000,
      )
    : null

  return (
    <div className="flex items-center gap-3 px-3 py-2 text-xs">
      {/* 任务信息 */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="font-medium truncate">{task.pluginId}</span>
          <span className="text-neutral-400">#{task.id.slice(0, 8)}</span>
          <span className={`inline-block h-2 w-2 shrink-0 rounded-full ${STATUS_COLORS[task.status] ?? ''} ${task.status === 'running' ? 'animate-pulse' : ''}`} />
        </div>
        <div className="text-neutral-500">
          {t(`task.status.${task.status}`)}
          {duration !== null && ` · ${duration}s`}
        </div>
      </div>

      {/* 操作 */}
      <div className="flex items-center gap-1 shrink-0">
        {task.status === 'running' && (
          <Button variant="ghost" size="sm" onClick={handleCancel} title={t('button.stop_task')}>
            {t('task.stop')}
          </Button>
        )}
        <Button variant="ghost" size="sm" onClick={handleOpenOutput} title={t('button.view_output')}>
          {t('panel.output')}
        </Button>
        {task.status !== 'running' && task.status !== 'created' && task.status !== 'queued' && (
          <Button variant="ghost" size="sm" onClick={handleDelete} title={t('button.delete_task')}>
            ✕
          </Button>
        )}
      </div>
    </div>
  )
}

interface TaskPanelProps {
  onTaskClick?: (taskId: string) => void
}

export function TaskPanel({ onTaskClick }: TaskPanelProps) {
  const { t } = useT()
  const { tasks } = useTaskStore()

  useEffect(() => {
    useTaskStore.getState().fetchTasks()
  }, [])

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between border-b border-border px-3 py-1.5">
        <span className="text-xs font-semibold uppercase text-neutral-500">
          {t('panel.tasks')} ({tasks.length})
        </span>
      </div>
      <div className="flex-1 overflow-y-auto">
        <div className="divide-y divide-neutral-100 dark:divide-neutral-800">
          {tasks.length === 0 ? (
            <div className="flex h-20 items-center justify-center text-xs text-neutral-400">
              {t('task.no_tasks')}
            </div>
          ) : (
            tasks.map((task) => <TaskItem key={task.id} task={task} onTaskClick={onTaskClick ?? (() => {})} />)
          )}
        </div>
      </div>
    </div>
  )
}
