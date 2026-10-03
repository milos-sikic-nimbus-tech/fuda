import type { ReactNode } from 'react'
import { ChartCard } from '@/components/organisms/insights/ChartCard'
import { Hint } from '@/components/atoms/Hint'
import type { Slice } from '@/lib/insights'

export function RankedBars({
  title,
  description,
  slices,
  labelOf,
  onSelect,
}: {
  title: string
  description?: string
  slices: Slice[]
  labelOf?: (key: string) => ReactNode
  onSelect?: (key: string) => void
}) {
  const max = Math.max(1, ...slices.map((s) => s.count))
  return (
    <ChartCard title={title} description={description}>
      {slices.length === 0 ? (
        <p className="py-6 text-center text-xs text-muted-foreground">Nothing to show</p>
      ) : (
        <ul className="space-y-1.5">
          {slices.map((s) => {
            const clickable = onSelect && s.key !== 'Other' && s.key !== 'None'
            const row = (
              <button
                type="button"
                disabled={!clickable}
                onClick={() => clickable && onSelect(s.key)}
                className="group grid w-full grid-cols-[minmax(0,9rem)_1fr_2.5rem] items-center gap-3 rounded-md px-1 py-0.5 text-left text-xs enabled:hover:bg-accent"
              >
                <span className="flex min-w-0 items-center gap-2 truncate text-foreground">
                  {labelOf ? labelOf(s.key) : s.key}
                </span>
                <span className="h-2.5 rounded-sm bg-muted">
                  <span
                    className="block h-full rounded-sm bg-series-1"
                    style={{ width: `${(s.count / max) * 100}%` }}
                  />
                </span>
                <span className="text-right text-muted-foreground tabular-nums">{s.count}</span>
              </button>
            )
            return (
              <li key={s.key}>
                {clickable ? <Hint label={`Open the board filtered to ${s.key}`}>{row}</Hint> : row}
              </li>
            )
          })}
        </ul>
      )}
    </ChartCard>
  )
}
