import { SlidersHorizontal } from 'lucide-react'
import { useMemo } from 'react'
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
import type { Card, Column, Facets, LabelGroup } from '@/lib/api'
import { countFacets, type FacetCounts, type ListKey, listOf } from '@/lib/filters'
import { cn } from '@/lib/utils'

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

function personOptions(names: string[], counts: Map<string, number>): Option[] {
  return names.map((name) => ({
    value: name,
    search: name,
    count: counts.get(name) ?? 0,
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

function useListFilter(key: ListKey) {
  const { search, toggle, update } = useBoardSearch()
  return {
    selected: listOf(search[key]),
    onToggle: (value: string) => toggle(key, value),
    onClear: () => update({ [key]: undefined }),
  }
}

function LabelGroupFilter({ group, counts }: { group: LabelGroup; counts: FacetCounts }) {
  const { search, toggle, update } = useBoardSearch()
  const prefix = `${group.name}:`
  const labels = listOf(search.label)
  return (
    <FilterSelect
      title={groupTitle(group.name)}
      options={group.values.map((v) => ({
        value: prefix + v,
        label: v,
        search: v,
        count: counts.label.get(prefix + v) ?? 0,
      }))}
      selected={labels.filter((l) => l.startsWith(prefix))}
      onToggle={(v) => toggle('label', v)}
      onClear={() =>
        update({ label: labels.filter((l) => !l.startsWith(prefix)).join(',') || undefined })
      }
    />
  )
}

export function FilterBar({
  facets,
  columns,
  cards,
  shown,
}: {
  facets: Facets
  columns: Column[]
  cards: Card[]
  shown: number
}) {
  const { search, update } = useBoardSearch()
  const counts = useMemo(() => countFacets(cards), [cards])
  const owner = useListFilter('owner')
  const tester = useListFilter('tester')
  const status = useListFilter('status')
  const prefix = useListFilter('prefix')

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
      search.text,
    ].some(Boolean)

  return (
    <div className="border-b border-border/60 bg-card/50 backdrop-blur-sm">
      <div className="flex flex-wrap items-center gap-2 px-4 py-2.5">
        <SearchInput value={search.q} onChange={(q) => update({ q })} />
        {primary.map((g) => (
          <LabelGroupFilter key={g.name} group={g} counts={counts} />
        ))}
        <FilterSelect
          title="Owner"
          options={personOptions(facets.people, counts.owner)}
          {...owner}
        />
        <Popover>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              size="sm"
              className={cn(
                'h-8 gap-1.5 bg-card shadow-xs',
                moreActive && 'border-primary/50 bg-primary/5 text-primary',
              )}
            >
              <SlidersHorizontal className="size-3.5" />
              More filters
            </Button>
          </PopoverTrigger>
          <PopoverContent align="start" className="w-80 space-y-4">
            <div className="flex flex-wrap gap-2 empty:hidden">
              {secondary.map((g) => (
                <LabelGroupFilter key={g.name} group={g} counts={counts} />
              ))}
              <FilterSelect
                title="Tester"
                options={personOptions(facets.testers, counts.tester)}
                {...tester}
              />
              <FilterSelect
                title="Status"
                options={(statusNeeded ? facets.statuses : []).map((s) => ({
                  value: s,
                  search: s,
                  count: counts.status.get(s) ?? 0,
                  label: (
                    <span className="inline-flex items-center gap-2">
                      <StatusDot status={s} />
                      {s}
                    </span>
                  ),
                }))}
                {...status}
              />
              <FilterSelect
                title="Prefix"
                options={(facets.prefixes.length > 1 ? facets.prefixes : []).map((p) => ({
                  value: p,
                  search: p,
                  count: counts.prefix.get(p) ?? 0,
                  label: <span className="font-mono">{p}</span>,
                }))}
                {...prefix}
              />
            </div>
            <DateRange
              from={search.from}
              to={search.to}
              onChange={(range) => update({ from: range.from, to: range.to })}
            />
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
            {shown === cards.length ? `${cards.length} tasks` : `${shown} of ${cards.length}`}
          </span>
          <Select
            value={search.sort ?? 'id'}
            onValueChange={(v) => update({ sort: v === 'added' ? 'added' : undefined })}
          >
            <SelectTrigger size="sm" className="h-8 w-40 bg-card text-xs shadow-xs">
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
