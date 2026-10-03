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
      label={
        <>
          {group ? `${group}: ${value}` : value}
          <span className="opacity-70"> · click to {active ? 'remove the' : ''} filter</span>
        </>
      }
    >
      <button
        type="button"
        onClick={(e) => {
          e.stopPropagation()
          onClick?.()
        }}
        className={cn(
          'inline-flex max-w-44 items-center gap-1 truncate rounded-md border px-1.5 py-0.5 text-[11px] leading-tight transition-colors focus-visible:ring-2 focus-visible:ring-ring/50 focus-visible:outline-none',
          active
            ? 'border-primary bg-primary/10 text-primary'
            : 'border-transparent bg-secondary text-secondary-foreground hover:border-primary/30',
        )}
      >
        {isType && (
          <span
            className="size-1.5 shrink-0 rounded-full"
            style={{ backgroundColor: typeHue(value) }}
          />
        )}
        {showGroup && group && !isType && group !== 'label' && (
          <span className="opacity-55">{group}</span>
        )}
        <span className="truncate">{value}</span>
      </button>
    </Hint>
  )
}
