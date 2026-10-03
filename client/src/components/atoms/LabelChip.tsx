import { cn } from '@/lib/utils'

export function labelText(label: string): { group: string; value: string } {
  const [group, ...rest] = label.split(':')
  return rest.length ? { group, value: rest.join(':') } : { group: '', value: group }
}

export function LabelChip({
  label,
  active,
  onClick,
}: {
  label: string
  active?: boolean
  onClick?: () => void
}) {
  const { group, value } = labelText(label)
  return (
    <button
      type="button"
      onClick={(e) => {
        e.stopPropagation()
        onClick?.()
      }}
      className={cn(
        'inline-flex max-w-40 items-center gap-1 truncate rounded-md border px-1.5 py-0.5 text-[11px] leading-none',
        active
          ? 'border-primary bg-primary/10 text-primary'
          : 'border-border bg-muted text-muted-foreground hover:border-primary/40',
      )}
    >
      {group && group !== 'label' && <span className="opacity-60">{group}</span>}
      <span className="truncate">{value}</span>
    </button>
  )
}
