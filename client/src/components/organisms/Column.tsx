import { CircleHelp } from 'lucide-react'
import { Hint } from '@/components/atoms/Hint'
import { StatusIcon } from '@/components/atoms/StatusIcon'
import { TaskCard } from '@/components/molecules/TaskCard'
import type { Card, Column as ColumnData } from '@/lib/api'

export function Column({
  column,
  cards,
  prLink,
}: {
  column: ColumnData
  cards: Card[]
  prLink?: string
}) {
  const status = column.prOpen ? 'in review' : (column.statuses[0] ?? '')
  return (
    <section className="flex w-76 shrink-0 flex-col rounded-xl bg-muted/50 dark:bg-muted/30">
      <header className="flex h-10 items-center gap-2 px-3 text-[13px] font-medium">
        <StatusIcon status={status} />
        <span>{column.name}</span>
        <span className="text-muted-foreground tabular-nums">{cards.length}</span>
        {column.unknown && (
          <Hint label="No stage in stages.md lists this status, so it gets its own column">
            <CircleHelp className="ml-auto size-3.5 text-amber-600" />
          </Hint>
        )}
      </header>
      <div className="flex min-h-0 flex-1 flex-col gap-1.5 overflow-y-auto px-1.5 pb-1.5">
        {cards.length === 0 ? (
          <p className="px-3 py-8 text-center text-xs text-muted-foreground/80">No tasks</p>
        ) : (
          cards.map((c) => <TaskCard key={c.id} card={c} column={column} prLink={prLink} />)
        )}
      </div>
    </section>
  )
}
