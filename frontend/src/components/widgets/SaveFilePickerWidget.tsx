import { useT } from '../../contexts/TranslationContext'
import { saveFileDialog } from '../../services/dialogApi'

export function SaveFilePickerWidget(props: any) {
  const { t } = useT()
  const { value, onChange } = props
  return (
    <div className="flex gap-2">
      <input
        type="text"
        className="flex-1 rounded-md border px-3 py-1.5 text-sm border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        value={value ?? ''}
        placeholder={t('form.save_file')}
        onChange={(e) => onChange(e.target.value)}
      />
      <button
        type="button"
        className="rounded-md border px-2 py-1 text-sm border-input bg-background hover:bg-accent"
        onClick={async () => {
          const file = await saveFileDialog({ title: t('form.save_file') })
          if (file) onChange(file)
        }}
      >
        ...
      </button>
    </div>
  )
}
