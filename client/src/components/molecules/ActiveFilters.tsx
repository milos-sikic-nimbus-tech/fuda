import { X } from 'lucide-react'
import type { ReactNode } from 'react'
import { labelText } from '@/components/atoms/LabelChip'
import { StatusIcon } from '@/components/atoms/StatusIcon'
import { Button } from '@/components/ui/button'
import { useBoardSearch } from '@/hooks/useBoardSearch'
import { typeHue } from '@/lib/colors'
import { type BoardSearch, type ListKey, listOf } from '@/lib/filters'

type Chip = {
  key: string
  group: string
  value: string
  marker?: ReactNode
  remove: Partial<BoardSearch> | [ListKey, string]
}

function chipsOf(search: BoardSearch): Chip[] {
  const chips: Chip[] = []
  const list = (key: ListKey, group: string, marker?: (v: string) => ReactNode) => {
    for (const value of listOf(search[key]))
      chips.push({
        key: `${key}:${value}`,
        group,
        value,
        marker: marker?.(value),
        remove: [key, value],
      })
  }
  list('status', 'status', (v) => <StatusIcon status={v} className="size-3" />)
  for (const label of listOf(search.label)) {
    const { group, value } = labelText(label)
    chips.push({
      key: `label:${label}`,
      group: group || 'label',
      value,
      marker:
        group === 'type' ? (
          <span className="size-1.5 rounded-full" style={{ backgroundColor: typeHue(value) }} />
        ) : undefined,
      remove: ['label', label],
    })
  }
  list('owner', 'owner')
  list('tester', 'tester')
  list('prefix', 'prefix')
  if (search.q)
    chips.push({
      key: 'q',
      group: search.text ? 'text' : 'search',
      value: search.q,
      remove: { q: undefined, text: undefined },
    })
  if (search.from || search.to)
    chips.push({
      key: 'added',
      group: 'added',
      value: `${search.from ?? '…'} – ${search.to ?? '…'}`,
      remove: { from: undefined, to: undefined },
    })
  if (search.blocked)
    chips.push({ key: 'blocked', group: 'only', value: 'blocked', remove: { blocked: undefined } })
  if (search.available)
    chips.push({
      key: 'available',
      group: 'only',
      value: 'available',
      remove: { available: undefined },
    })
  return chips
}

export function ActiveFilters() {
  const { search, update, toggle, clear } = useBoardSearch()
  const chips = chipsOf(search)
  if (chips.length === 0) return null
  return (
    <div className="flex flex-wrap items-center gap-1.5 px-3 pb-2.5 sm:px-4">
      {chips.map((chip) => (
        <button
          key={chip.key}
          type="button"
          aria-label={`Remove filter ${chip.group} ${chip.value}`}
          onClick={() =>
            Array.isArray(chip.remove) ? toggle(...chip.remove) : update(chip.remove)
          }
          className="group/chip inline-flex h-6 items-center gap-1.5 rounded-full border border-primary/25 bg-primary/8 pr-1.5 pl-2.5 text-xs transition-colors hover:border-primary/50 hover:bg-primary/12"
        >
          <span className="text-muted-foreground">{chip.group}</span>
          {chip.marker}
          <span className="max-w-48 truncate font-medium text-foreground">{chip.value}</span>
          <X className="size-3.5 rounded-full text-muted-foreground opacity-50 transition-opacity group-hover/chip:opacity-100" />
        </button>
      ))}
      <Button
        variant="ghost"
        size="sm"
        className="h-6 px-2 text-xs text-muted-foreground"
        onClick={clear}
      >
        Clear all
      </Button>
    </div>
  )
}
