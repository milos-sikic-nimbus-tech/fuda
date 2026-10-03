import { Hint } from '@/components/atoms/Hint'
import { typeHue } from '@/lib/colors'
import { cn } from '@/lib/utils'

export function labelText(label: string): { group: string; value: string } {
  const [group, ...rest] = label.split(':')
  return rest.length ? { group, value: rest.join(':') } : { group: '', value: group }
}

export function LabelChip({
  label,
  active,
  showGroup = true,
  onClick,
}: {
  label: string
  active?: boolean
  showGroup?: boolean
  onClick?: () => void
}) {
  const { group, value } = labelText(label)
  const isType = group === 'type'
  return (
    <Hint
      label={`${group ? `${group}: ` : ''}${value} · click to ${active ? 'remove the filter' : 'filter'}`}
    >
      <button
        type="button"
        onClick={(e) => {
          e.stopPropagation()
          onClick?.()
        }}
        className={cn(
          'inline-flex h-5 max-w-40 items-center gap-1.5 truncate rounded-full border px-2 text-[11px] transition-colors',
          'focus-visible:ring-2 focus-visible:ring-ring/50 focus-visible:outline-none',
          active
            ? 'border-primary/40 bg-primary/10 text-primary'
            : 'border-border text-muted-foreground hover:border-foreground/20 hover:text-foreground',
        )}
      >
        {isType && (
          <span
            className="size-1.5 shrink-0 rounded-full"
            style={{ backgroundColor: typeHue(value) }}
          />
        )}
        {showGroup && group && !isType && group !== 'label' && (
          <span className="opacity-60">{group}</span>
        )}
        <span className="truncate">{value}</span>
      </button>
    </Hint>
  )
}
