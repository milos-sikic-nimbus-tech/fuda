import { Column } from '@/components/organisms/Column'
import { useCollapsedColumns } from '@/hooks/useCollapsedColumns'
import type { Card, Column as ColumnData } from '@/lib/api'

export function Board({
  columns,
  cards,
  prLink,
}: {
  columns: ColumnData[]
  cards: Card[]
  prLink?: string
}) {
  const { isCollapsed, toggle } = useCollapsedColumns()
  return (
    <div className="flex min-h-0 flex-1 snap-x snap-mandatory gap-3 overflow-x-auto px-3 pb-3 sm:snap-none sm:px-4 sm:pb-4">
      {columns.map((column) => (
        <Column
          key={column.id}
          column={column}
          cards={cards.filter((c) => c.column === column.id)}
          prLink={prLink}
          collapsed={isCollapsed(column.id)}
          onToggleCollapse={() => toggle(column.id)}
        />
      ))}
    </div>
  )
}
