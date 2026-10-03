import { CircleHelp } from 'lucide-react'
import { StatusDot } from '@/components/atoms/StatusDot'
import { TaskCard } from '@/components/molecules/TaskCard'
import type { Card, Column as ColumnData } from '@/lib/api'

export function Column({ column, cards }: { column: ColumnData; cards: Card[] }) {
  return (
    <section className="flex w-72 shrink-0 flex-col rounded-xl bg-muted/60">
      <header className="flex items-center gap-2 px-3 pt-3 pb-2 text-sm font-medium">
        <StatusDot status={column.prOpen ? 'in review' : (column.statuses[0] ?? '')} />
        {column.name}
        {column.unknown && (
          <span
            title="This status is not in stages.md"
            className="inline-flex items-center gap-1 text-xs font-normal text-amber-700"
          >
            <CircleHelp className="size-3.5" />
            unknown
          </span>
        )}
        <span className="ml-auto text-xs font-normal text-muted-foreground">{cards.length}</span>
      </header>
      <div className="flex min-h-0 flex-1 flex-col gap-2 overflow-y-auto px-2 pb-2">
        {cards.map((c) => (
          <TaskCard key={c.id} card={c} column={column} />
        ))}
      </div>
    </section>
  )
}
