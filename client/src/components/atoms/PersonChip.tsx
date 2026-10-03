import { Hint } from '@/components/atoms/Hint'
import { initials, personColor } from '@/lib/colors'
import { cn } from '@/lib/utils'

export function Avatar({
  name,
  hint,
  className,
}: {
  name: string
  hint?: string
  className?: string
}) {
  const avatar = (
    <span
      style={{ backgroundColor: personColor(name) }}
      className={cn(
        'inline-flex size-5 shrink-0 items-center justify-center rounded-full text-[9px] font-semibold text-white',
        className,
      )}
    >
      {initials(name)}
    </span>
  )
  return hint ? <Hint label={hint}>{avatar}</Hint> : avatar
}

export function PersonChip({ name, onClick }: { name: string; onClick?: () => void }) {
  return (
    <button
      type="button"
      onClick={(e) => {
        e.stopPropagation()
        onClick?.()
      }}
      className="inline-flex items-center gap-1.5 rounded-full border border-border py-0.5 pr-2 pl-0.5 text-xs hover:border-primary/40"
    >
      <Avatar name={name} />
      {name}
    </button>
  )
}
