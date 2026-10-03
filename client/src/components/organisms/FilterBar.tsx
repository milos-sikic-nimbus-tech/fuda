import { SlidersHorizontal } from 'lucide-react'
import { Avatar } from '@/components/atoms/PersonChip'
import { StatusDot } from '@/components/atoms/StatusDot'
import { ActiveFilters } from '@/components/molecules/ActiveFilters'
import { DateRange } from '@/components/molecules/DateRange'
import { FilterSelect, type Option } from '@/components/molecules/FilterSelect'
import { SearchInput } from '@/components/molecules/SearchInput'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useBoardSearch } from '@/hooks/useBoardSearch'
import type { Column, Facets, LabelGroup } from '@/lib/api'
import { listOf } from '@/lib/filters'

const groupTitles: Record<string, string> = {
  type: 'Type',
  epic: 'Epic',
  area: 'Area',
  label: 'Labels',
}
const primaryGroups = ['type', 'epic']

function groupTitle(group: string): string {
  return groupTitles[group] ?? group.charAt(0).toUpperCase() + group.slice(1)
}

function personOptions(names: string[]): Option[] {
  return names.map((name) => ({
    value: name,
    search: name,
    label: (
      <span className="inline-flex items-center gap-2">
        <Avatar name={name} />
        {name}
      </span>
    ),
  }))
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
    <label className="flex cursor-pointer items-center gap-2 text-sm select-none">
      <Checkbox checked={checked} onCheckedChange={(v) => onChange(v === true)} />
      {label}
    </label>
  )
}

function LabelGroupFilter({ group }: { group: LabelGroup }) {
  const { search, toggle } = useBoardSearch()
  return (
    <FilterSelect
      title={groupTitle(group.name)}
      options={group.values.map((v) => ({ value: `${group.name}:${v}`, label: v, search: v }))}
      selected={listOf(search.label).filter((l) => l.startsWith(`${group.name}:`))}
      onToggle={(v) => toggle('label', v)}
    />
  )
}

export function FilterBar({
  facets,
  columns,
  shown,
  total,
}: {
  facets: Facets
  columns: Column[]
  shown: number
  total: number
}) {
  const { search, update, toggle } = useBoardSearch()
  const statusNeeded = columns.some((c) => c.statuses.length > 1 || c.unknown)
  const primary = facets.labelGroups.filter((g) => primaryGroups.includes(g.name))
  const secondary = facets.labelGroups.filter((g) => !primaryGroups.includes(g.name))
  const moreActive =
    secondary.some((g) => listOf(search.label).some((l) => l.startsWith(`${g.name}:`))) ||
    [
      search.tester,
      search.status,
      search.prefix,
      search.from,
      search.to,
      search.blocked,
      search.available,
    ].some(Boolean)

  return (
    <div className="border-b border-border/60">
      <div className="flex flex-wrap items-center gap-2 px-4 py-3">
        <SearchInput value={search.q} onChange={(q) => update({ q })} />
        {primary.map((g) => (
          <LabelGroupFilter key={g.name} group={g} />
        ))}
        <FilterSelect
          title="Owner"
          options={personOptions(facets.people)}
          selected={listOf(search.owner)}
          onToggle={(v) => toggle('owner', v)}
        />
        <Popover>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              size="sm"
              className="h-8 gap-1.5"
              data-active={moreActive || undefined}
            >
              <SlidersHorizontal className="size-3.5" />
              More filters
              {moreActive && <span className="size-1.5 rounded-full bg-primary" />}
            </Button>
          </PopoverTrigger>
          <PopoverContent align="start" className="w-80 space-y-4">
            {secondary.length +
              (facets.testers.length ? 1 : 0) +
              (statusNeeded ? 1 : 0) +
              (facets.prefixes.length > 1 ? 1 : 0) >
              0 && (
              <div className="flex flex-wrap gap-2">
                {secondary.map((g) => (
                  <LabelGroupFilter key={g.name} group={g} />
                ))}
                <FilterSelect
                  title="Tester"
                  options={personOptions(facets.testers)}
                  selected={listOf(search.tester)}
                  onToggle={(v) => toggle('tester', v)}
                />
                <FilterSelect
                  title="Status"
                  options={(statusNeeded ? facets.statuses : []).map((s) => ({
                    value: s,
                    search: s,
                    label: (
                      <span className="inline-flex items-center gap-2">
                        <StatusDot status={s} />
                        {s}
                      </span>
                    ),
                  }))}
                  selected={listOf(search.status)}
                  onToggle={(v) => toggle('status', v)}
                />
                <FilterSelect
                  title="Prefix"
                  options={(facets.prefixes.length > 1 ? facets.prefixes : []).map((p) => ({
                    value: p,
                    search: p,
                    label: <span className="font-mono">{p}</span>,
                  }))}
                  selected={listOf(search.prefix)}
                  onToggle={(v) => toggle('prefix', v)}
                />
              </div>
            )}
            <DateRange from={search.from} to={search.to} onChange={(range) => update(range)} />
            <div className="space-y-2">
              <Toggle
                label="Blocked only"
                checked={!!search.blocked}
                onChange={(v) => update({ blocked: v || undefined })}
              />
              <Toggle
                label="Available to pick up"
                checked={!!search.available}
                onChange={(v) => update({ available: v || undefined })}
              />
              <Toggle
                label="Search also in task text"
                checked={!!search.text}
                onChange={(v) => update({ text: v || undefined })}
              />
            </div>
          </PopoverContent>
        </Popover>
        <div className="ml-auto flex items-center gap-3">
          <span className="text-xs text-muted-foreground tabular-nums">
            {shown === total ? `${total} tasks` : `${shown} of ${total}`}
          </span>
          <Select
            value={search.sort ?? 'id'}
            onValueChange={(v) => update({ sort: v === 'added' ? 'added' : undefined })}
          >
            <SelectTrigger size="sm" className="h-8 w-40 text-xs">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="id">Sort by id</SelectItem>
              <SelectItem value="added">Newest added first</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>
      <ActiveFilters />
    </div>
  )
}
