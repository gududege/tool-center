import { useT } from '../../contexts/TranslationContext'
import { openDirectoryDialog } from '../../services/dialogApi'

export function FolderPickerWidget(props: any) {
  const { t } = useT()
  const { value, onChange } = props
  return (
    <div className="flex gap-2">
      <input
        type="text"
        className="flex-1 rounded-md border px-3 py-1.5 text-sm border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        value={value ?? ''}
        placeholder={t('form.select_directory')}
        onChange={(e) => onChange(e.target.value)}
      />
      <button
        type="button"
        className="rounded-md border px-2 py-1 text-sm border-input bg-background hover:bg-accent"
        onClick={async () => {
          const dir = await openDirectoryDialog({ title: t('form.select_directory') })
          if (dir) onChange(dir)
        }}
      >
        ...
      </button>
    </div>
  )
}
