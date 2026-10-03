import { X } from 'lucide-react'
import { labelText } from '@/components/atoms/LabelChip'
import { Button } from '@/components/ui/button'
import { useBoardSearch } from '@/hooks/useBoardSearch'
import { type BoardSearch, type ListKey, listOf } from '@/lib/filters'

type Chip = { key: string; text: string; remove: Partial<BoardSearch> | [ListKey, string] }

function chipsOf(search: BoardSearch): Chip[] {
  const lists: [ListKey, string][] = [
    ['status', 'status'],
    ['owner', 'owner'],
    ['tester', 'tester'],
    ['prefix', 'prefix'],
  ]
  const chips: Chip[] = lists.flatMap(([key, name]) =>
    listOf(search[key]).map((v) => ({
      key: `${key}:${v}`,
      text: `${name}: ${v}`,
      remove: [key, v] as [ListKey, string],
    })),
  )
  for (const l of listOf(search.label)) {
    const { group, value } = labelText(l)
    chips.push({
      key: `label:${l}`,
      text: group ? `${group}: ${value}` : value,
      remove: ['label', l],
    })
  }
  if (search.q)
    chips.push({
      key: 'q',
      text: `“${search.q}”${search.text ? ' in text' : ''}`,
      remove: { q: undefined, text: undefined },
    })
  if (search.from || search.to)
    chips.push({
      key: 'added',
      text: `added ${search.from ?? '…'} – ${search.to ?? '…'}`,
      remove: { from: undefined, to: undefined },
    })
  if (search.blocked)
    chips.push({ key: 'blocked', text: 'blocked', remove: { blocked: undefined } })
  if (search.available)
    chips.push({ key: 'available', text: 'available', remove: { available: undefined } })
  return chips
}

export function ActiveFilters() {
  const { search, update, toggle, clear } = useBoardSearch()
  const chips = chipsOf(search)
  if (chips.length === 0) return null
  return (
    <div className="flex flex-wrap items-center gap-1.5 px-4 pb-2">
      {chips.map((chip) => (
        <button
          key={chip.key}
          type="button"
          onClick={() =>
            Array.isArray(chip.remove) ? toggle(...chip.remove) : update(chip.remove)
          }
          className="inline-flex items-center gap-1 rounded-full border border-primary/30 bg-primary/8 px-2.5 py-0.5 text-xs text-primary hover:bg-primary/15"
        >
          {chip.text}
          <X className="size-3" />
        </button>
      ))}
      <Button variant="ghost" size="sm" className="h-6 px-2 text-xs" onClick={clear}>
        Clear all
      </Button>
    </div>
  )
}
