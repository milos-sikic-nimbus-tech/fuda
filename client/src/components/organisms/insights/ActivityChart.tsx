import { useState } from 'react'
import { ChartCard, LegendItem } from '@/components/organisms/insights/ChartCard'
import { useElementWidth } from '@/hooks/useElementWidth'
import type { Bucket, BucketUnit } from '@/lib/insights'

const height = 200
const padding = { top: 12, right: 8, bottom: 24, left: 32 }
const series = [
  { key: 'created', label: 'Created', color: 'var(--color-series-1)' },
  { key: 'claimed', label: 'Claimed', color: 'var(--color-series-2)' },
] as const

function niceMax(value: number): number {
  if (value <= 4) return 4
  const step = 10 ** Math.floor(Math.log10(value))
  return Math.ceil(value / step) * step
}

export function ActivityChart({ unit, buckets }: { unit: BucketUnit; buckets: Bucket[] }) {
  const { ref, width } = useElementWidth<HTMLDivElement>()
  const [hover, setHover] = useState<number | null>(null)
  const max = niceMax(Math.max(1, ...buckets.flatMap((b) => [b.created, b.claimed])))
  const plotWidth = Math.max(0, width - padding.left - padding.right)
  const plotHeight = height - padding.top - padding.bottom
  const band = buckets.length ? plotWidth / buckets.length : 0
  const barWidth = Math.max(2, Math.min(18, (band - 6) / 2 - 1))
  const y = (v: number) => padding.top + plotHeight - (v / max) * plotHeight
  const ticks = [0, max / 2, max]
  const labelEvery = Math.ceil(buckets.length / Math.max(1, Math.floor(plotWidth / 56)))
  const active = hover === null ? null : buckets[hover]
  const periodName = unit === 'day' ? 'Day' : 'Week of'

  return (
    <ChartCard
      title={`Created vs claimed, per ${unit}`}
      description="From each task's added and claimed dates."
      legend={
        <div className="flex items-center gap-3">
          {series.map((s) => (
            <LegendItem key={s.key} color={s.color} label={s.label} />
          ))}
        </div>
      }
      table={
        <div className="max-h-56 overflow-y-auto text-xs">
          <table className="w-full">
            <thead className="text-muted-foreground">
              <tr>
                <th className="py-1 text-left font-medium">{periodName}</th>
                <th className="py-1 text-right font-medium">Created</th>
                <th className="py-1 text-right font-medium">Claimed</th>
              </tr>
            </thead>
            <tbody className="tabular-nums">
              {buckets.map((w) => (
                <tr key={w.start} className="border-t border-border">
                  <td className="py-1">{w.label}</td>
                  <td className="py-1 text-right">{w.created}</td>
                  <td className="py-1 text-right">{w.claimed}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      }
    >
      <div ref={ref} className="relative" onMouseLeave={() => setHover(null)}>
        {width > 0 && (
          <svg
            width={width}
            height={height}
            role="img"
            aria-label={`Tasks created and claimed per ${unit}`}
          >
            {ticks.map((t) => (
              <g key={t}>
                <line
                  x1={padding.left}
                  x2={width - padding.right}
                  y1={y(t)}
                  y2={y(t)}
                  className="stroke-border"
                  strokeDasharray={t === 0 ? undefined : '2 3'}
                />
                <text
                  x={padding.left - 6}
                  y={y(t)}
                  dy="0.32em"
                  textAnchor="end"
                  className="fill-muted-foreground text-[10px] tabular-nums"
                >
                  {t}
                </text>
              </g>
            ))}
            {buckets.map((w, i) => {
              const x0 = padding.left + i * band + band / 2 - barWidth - 1
              return (
                <g key={w.start} onMouseEnter={() => setHover(i)}>
                  <rect
                    x={padding.left + i * band}
                    y={padding.top}
                    width={band}
                    height={plotHeight}
                    fill={hover === i ? 'var(--color-muted)' : 'transparent'}
                  />
                  {series.map((s, j) => {
                    const value = w[s.key]
                    const top = y(value)
                    const h = padding.top + plotHeight - top
                    return value > 0 ? (
                      <path
                        key={s.key}
                        d={`M${x0 + j * (barWidth + 2)},${top + h} v${-Math.max(0, h - 3)} q0,-3 3,-3 h${barWidth - 6} q3,0 3,3 v${Math.max(0, h - 3)} z`}
                        fill={s.color}
                      />
                    ) : null
                  })}
                  {i % labelEvery === 0 && (
                    <text
                      x={padding.left + i * band + band / 2}
                      y={height - 6}
                      textAnchor="middle"
                      className="fill-muted-foreground text-[10px]"
                    >
                      {w.label}
                    </text>
                  )}
                </g>
              )
            })}
          </svg>
        )}
        {active && hover !== null && (
          <div
            className="pointer-events-none absolute top-1 z-10 rounded-md border border-border bg-popover px-2.5 py-1.5 text-xs shadow-md"
            style={{
              left: Math.min(width - 140, Math.max(0, padding.left + hover * band + band / 2 - 70)),
            }}
          >
            <p className="mb-1 font-medium">
              {unit === 'day' ? active.label : `Week of ${active.label}`}
            </p>
            {series.map((s) => (
              <p
                key={s.key}
                className="flex items-center justify-between gap-4 text-muted-foreground"
              >
                <span className="inline-flex items-center gap-1.5">
                  <span className="size-2 rounded-sm" style={{ backgroundColor: s.color }} />
                  {s.label}
                </span>
                <span className="text-foreground tabular-nums">{active[s.key]}</span>
              </p>
            ))}
          </div>
        )}
      </div>
    </ChartCard>
  )
}
