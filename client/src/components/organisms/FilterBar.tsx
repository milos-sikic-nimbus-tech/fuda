import { LayoutGrid, List } from 'lucide-react'
import { Hint } from '@/components/atoms/Hint'
import { ToggleGroup, ToggleGroupItem } from '@/components/ui/toggle-group'
import { ActiveFilters } from '@/components/molecules/ActiveFilters'
import { FilterMenu } from '@/components/molecules/FilterMenu'
import { SearchInput } from '@/components/molecules/SearchInput'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useBoardSearch } from '@/hooks/useBoardSearch'
import { useFilterGroups } from '@/hooks/useFilterGroups'
import type { Card, Column, Facets } from '@/lib/api'

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
  const groups = useFilterGroups(facets, columns, cards)

  return (
    <div className="flex flex-wrap items-center gap-2 border-b border-border/60 bg-card/50 px-3 py-2.5 backdrop-blur-sm sm:px-4">
      <SearchInput value={search.q} onChange={(q) => update({ q })} />
      <FilterMenu groups={groups} />
      <ActiveFilters groups={groups} />
      <div className="ml-auto flex items-center gap-3">
        <ToggleGroup
          type="single"
          size="sm"
          variant="outline"
          value={search.view ?? 'board'}
          onValueChange={(v) => v && update({ view: v === 'list' ? 'list' : undefined })}
          className="bg-card shadow-xs"
        >
          <Hint label="Board">
            <ToggleGroupItem value="board" aria-label="Board view" className="h-8 px-2">
              <LayoutGrid className="size-3.5" />
            </ToggleGroupItem>
          </Hint>
          <Hint label="List">
            <ToggleGroupItem value="list" aria-label="List view" className="h-8 px-2">
              <List className="size-3.5" />
            </ToggleGroupItem>
          </Hint>
        </ToggleGroup>
        <span className="text-xs text-muted-foreground tabular-nums">
          {shown === cards.length ? `${cards.length} tasks` : `${shown} of ${cards.length}`}
        </span>
        <Select
          value={search.sort ?? 'id'}
          onValueChange={(v) => update({ sort: v === 'added' ? 'added' : undefined })}
        >
          <SelectTrigger size="sm" className="h-8 w-36 bg-card text-xs shadow-xs">
            <SelectValue />
          </SelectTrigger>
          <SelectContent align="end">
            <SelectItem value="id">Sort by id</SelectItem>
            <SelectItem value="added">Newest first</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  )
}
