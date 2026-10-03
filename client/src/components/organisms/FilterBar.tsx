import { X } from 'lucide-react'
import { Avatar } from '@/components/atoms/PersonChip'
import { StatusDot } from '@/components/atoms/StatusDot'
import { DateRange } from '@/components/molecules/DateRange'
import { FilterSelect } from '@/components/molecules/FilterSelect'
import { SearchInput } from '@/components/molecules/SearchInput'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { useBoardSearch } from '@/hooks/useBoardSearch'
import type { Facets } from '@/lib/api'
import { activeFilterCount, listOf } from '@/lib/filters'

const groupTitles: Record<string, string> = {
  type: 'Type',
  epic: 'Epic',
  area: 'Area',
  label: 'Labels',
}

function title(group: string): string {
  return groupTitles[group] ?? group.charAt(0).toUpperCase() + group.slice(1)
}

function Toggle({
  label,
  checked,
  onChange,
}: {
  label: string
  checked: boolean
  onChange: (v: boolean) => void
}) {
  return (
    <label className="flex cursor-pointer items-center gap-1.5 text-xs text-muted-foreground select-none">
      <Checkbox checked={checked} onCheckedChange={(v) => onChange(v === true)} />
      {label}
    </label>
  )
}

export function FilterBar({
  facets,
  shown,
  total,
}: {
  facets: Facets
  shown: number
  total: number
}) {
  const { search, update, toggle, clear } = useBoardSearch()
  const labels = listOf(search.label)
  const people = facets.people.map((p) => ({
    value: p,
    label: (
      <span className="inline-flex items-center gap-2">
        <Avatar name={p} />
        {p}
      </span>
    ),
    search: p,
  }))

  return (
    <div className="flex flex-wrap items-center gap-2 px-4 py-3">
      <SearchInput value={search.q} onChange={(q) => update({ q })} />
      <Toggle
        label="also in task text"
        checked={!!search.text}
        onChange={(v) => update({ text: v || undefined })}
      />
      <FilterSelect
        title="Status"
        options={facets.statuses.map((s) => ({
          value: s,
          label: (
            <span className="inline-flex items-center gap-2">
              <StatusDot status={s} />
              {s}
            </span>
          ),
          search: s,
        }))}
        selected={listOf(search.status)}
        onToggle={(v) => toggle('status', v)}
      />
      <FilterSelect
        title="Owner"
        options={people}
        selected={listOf(search.owner)}
        onToggle={(v) => toggle('owner', v)}
      />
      <FilterSelect
        title="Tester"
        options={people}
        selected={listOf(search.tester)}
        onToggle={(v) => toggle('tester', v)}
      />
      {facets.labelGroups.map((g) => (
        <FilterSelect
          key={g.name}
          title={title(g.name)}
          options={g.values.map((v) => ({ value: `${g.name}:${v}`, label: v, search: v }))}
          selected={labels.filter((l) => l.startsWith(`${g.name}:`))}
          onToggle={(v) => toggle('label', v)}
        />
      ))}
      {facets.prefixes.length > 1 && (
        <FilterSelect
          title="Prefix"
          options={facets.prefixes.map((p) => ({
            value: p,
            label: <span className="font-mono">{p}</span>,
            search: p,
          }))}
          selected={listOf(search.prefix)}
          onToggle={(v) => toggle('prefix', v)}
        />
      )}
      <DateRange from={search.from} to={search.to} onChange={(range) => update(range)} />
      <Toggle
        label="Blocked"
        checked={!!search.blocked}
        onChange={(v) => update({ blocked: v || undefined })}
      />
      <Toggle
        label="Available"
        checked={!!search.available}
        onChange={(v) => update({ available: v || undefined })}
      />
      <select
        value={search.sort ?? 'id'}
        onChange={(e) => update({ sort: e.target.value === 'added' ? 'added' : undefined })}
        className="h-8 rounded-md border border-input bg-card px-2 text-xs"
      >
        <option value="id">Sort: id</option>
        <option value="added">Sort: newest added</option>
      </select>
      <span className="ml-auto text-xs text-muted-foreground">
        {shown === total ? `${total} tasks` : `${shown} of ${total} tasks`}
      </span>
      {activeFilterCount(search) > 0 && (
        <Button variant="ghost" size="sm" className="h-8 gap-1" onClick={clear}>
          <X className="size-3.5" />
          Clear
        </Button>
      )}
    </div>
  )
}
