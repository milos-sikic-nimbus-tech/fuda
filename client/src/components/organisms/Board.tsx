import { Column } from '@/components/organisms/Column'
import type { Card, Column as ColumnData } from '@/lib/api'

export function Board({ columns, cards }: { columns: ColumnData[]; cards: Card[] }) {
  return (
    <div className="flex min-h-0 flex-1 gap-3 overflow-x-auto px-4 pb-4">
      {columns.map((column) => (
        <Column
          key={column.id}
          column={column}
          cards={cards.filter((c) => c.column === column.id)}
        />
      ))}
    </div>
  )
}
