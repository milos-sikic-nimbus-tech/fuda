import { X } from 'lucide-react'
import type { ReactNode } from 'react'
import { ValueList } from '@/components/molecules/ValueList'
import { Command } from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import type { FilterGroup } from '@/hooks/useFilterGroups'
import { cn } from '@/lib/utils'

const segment = 'flex h-full items-center px-2 transition-colors hover:bg-accent'

export function PillShell({
  children,
  onRemove,
  removeLabel,
}: {
  children: ReactNode
  onRemove: () => void
  removeLabel: string
}) {
  return (
    <span className="inline-flex h-7 items-stretch overflow-hidden rounded-md border border-border bg-card text-xs shadow-xs">
      {children}
      <button
        type="button"
        aria-label={removeLabel}
        onClick={onRemove}
        className={cn(segment, 'border-l text-muted-foreground hover:text-foreground')}
      >
        <X className="size-3.5" />
      </button>
    </span>
  )
}

function valueSummary(group: FilterGroup): ReactNode {
  const chosen = group.values.filter((v) => group.selected.includes(v.value))
  if (chosen.length === 1) return chosen[0].label
  if (chosen.length <= 2) return chosen.map((v) => v.text).join(', ')
  return `${chosen.length} ${group.plural}`
}

export function FilterPill({ group }: { group: FilterGroup }) {
  return (
    <PillShell onRemove={group.clear} removeLabel={`Remove ${group.title} filter`}>
      <span className="flex items-center px-2 font-medium">{group.title}</span>
      <span className="flex items-center border-l px-2 text-muted-foreground">
        {group.selected.length > 1 ? 'is any of' : 'is'}
      </span>
      <Popover>
        <PopoverTrigger asChild>
          <button type="button" className={cn(segment, 'max-w-56 border-l font-medium')}>
            <span className="truncate">{valueSummary(group)}</span>
          </button>
        </PopoverTrigger>
        <PopoverContent align="start" className="w-[min(18rem,calc(100vw-1.5rem))] p-0">
          <Command>
            <ValueList group={group} autoFocus />
          </Command>
        </PopoverContent>
      </Popover>
    </PillShell>
  )
}
