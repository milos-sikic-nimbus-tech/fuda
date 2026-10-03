import { createFileRoute } from '@tanstack/react-router'
import { useMemo } from 'react'
import { Board } from '@/components/organisms/Board'
import { FilterBar } from '@/components/organisms/FilterBar'
import { TaskSheet } from '@/components/organisms/TaskSheet'
import { boardSearchSchema, filterCards, sortCards } from '@/lib/filters'
import { useBoard, useTextSearch } from '@/lib/queries'

export const Route = createFileRoute('/')({
  validateSearch: boardSearchSchema,
  component: BoardPage,
})

function BoardPage() {
  const search = Route.useSearch()
  const { data, error } = useBoard()
  const { data: textMatches } = useTextSearch(search.q, !!search.text)

  const cards = useMemo(
    () => (data ? sortCards(filterCards(data.cards, search, textMatches), search.sort) : []),
    [data, search, textMatches],
  )
  const taskIds = useMemo(() => data?.cards.map((c) => c.id) ?? [], [data])

  if (error) {
    return (
      <p className="p-6 text-sm text-muted-foreground">
        No board data yet: {error.message}. The first sync may still be running.
      </p>
    )
  }
  if (!data) return <p className="p-6 text-sm text-muted-foreground">Loading the board…</p>

  return (
    <main className="flex min-h-0 flex-1 flex-col">
      <FilterBar facets={data.facets} shown={cards.length} total={data.cards.length} />
      <Board columns={data.columns} cards={cards} prLink={data.prLink} />
      <TaskSheet taskIds={taskIds} prLink={data.prLink} />
    </main>
  )
}
