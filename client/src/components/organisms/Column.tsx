import { CircleHelp } from 'lucide-react'
import { Hint } from '@/components/atoms/Hint'
import { TaskCard } from '@/components/molecules/TaskCard'
import type { Card, Column as ColumnData } from '@/lib/api'
import { hueStyle, statusHue } from '@/lib/colors'

export function Column({
  column,
  cards,
  prLink,
}: {
  column: ColumnData
  cards: Card[]
  prLink?: string
}) {
  const hue = statusHue(column.prOpen ? 'in review' : (column.statuses[0] ?? ''))
  return (
    <section
      style={hueStyle(hue)}
      className="flex w-72 shrink-0 flex-col rounded-xl border tint-border tint-surface"
    >
      <header className="flex items-center gap-2 px-3 pt-3 pb-2 text-[13px] font-semibold">
        <span className="h-3.5 w-1 rounded-full" style={{ backgroundColor: hue }} />
        {column.name}
        {column.unknown && (
          <Hint label="No stage in stages.md lists this status, so it gets its own column">
            <span className="inline-flex items-center gap-1 text-xs font-normal text-amber-700">
              <CircleHelp className="size-3.5" />
              unknown
            </span>
          </Hint>
        )}
        <span className="ml-auto rounded-full bg-card/80 px-2 py-0.5 text-[11px] font-medium text-muted-foreground tabular-nums">
          {cards.length}
        </span>
      </header>
      <div className="flex min-h-0 flex-1 flex-col gap-2 overflow-y-auto px-2 pb-2">
        {cards.length === 0 ? (
          <p className="mx-1 mt-1 rounded-lg border border-dashed tint-border px-3 py-6 text-center text-xs text-muted-foreground">
            Nothing here
          </p>
        ) : (
          cards.map((c) => <TaskCard key={c.id} card={c} column={column} prLink={prLink} />)
        )}
      </div>
    </section>
  )
}
