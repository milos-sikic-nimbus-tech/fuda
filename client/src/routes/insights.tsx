import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { parseISO, subWeeks } from 'date-fns'
import { useMemo, useState } from 'react'
import { Avatar } from '@/components/atoms/PersonChip'
import { StatusIcon } from '@/components/atoms/StatusIcon'
import { StatTile } from '@/components/molecules/StatTile'
import { RankedBars } from '@/components/organisms/insights/RankedBars'
import { WeeklyChart } from '@/components/organisms/insights/WeeklyChart'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useLabelColor } from '@/hooks/useLabelColor'
import {
  countBy,
  insightsSearchSchema,
  isOpen,
  labelsOf,
  type RangeKey,
  ranges,
  topSlices,
  weeklyActivity,
} from '@/lib/insights'
import { useArchive, useBoard } from '@/lib/queries'

export const Route = createFileRoute('/insights')({
  validateSearch: insightsSearchSchema,
  component: InsightsPage,
})

function InsightsPage() {
  const { range = '12w' } = Route.useSearch()
  const navigate = useNavigate({ from: '/insights' })
  const toBoard = useNavigate()
  const { data } = useBoard()
  const { data: archive = [] } = useArchive()
  const colorOf = useLabelColor()
  const weeks = ranges[range]
  const [today] = useState(() => new Date())

  const view = useMemo(() => {
    const cards = data?.cards ?? []
    const since = subWeeks(today, weeks)
    const within = (date: string) => date !== '' && parseISO(date) >= since
    const open = cards.filter(isOpen)
    return {
      weekly: weeklyActivity(cards, weeks, today),
      created: cards.filter((c) => within(c.added)).length,
      claimed: cards.filter((c) => within(c.claimed)).length,
      done: archive.filter((c) => within(c.done ?? '')).length,
      open: open.length,
      inProgress: cards.filter((c) => c.status.trim().toLowerCase() === 'in progress').length,
      blocked: cards.filter((c) => c.blockedBy).length,
      byStage: (data?.columns ?? []).map((col) => ({
        key: col.name,
        count: cards.filter((c) => c.column === col.id).length,
        status: col.prOpen ? 'in review' : (col.statuses[0] ?? ''),
      })),
      byEpic: topSlices(countBy(open, labelsOf('epic'))),
      byType: topSlices(countBy(open, labelsOf('type'))),
      byOwner: topSlices(
        countBy(
          cards.filter((c) => c.owners.length && isOpen(c)),
          (c) => c.owners,
        ),
      ),
    }
  }, [data, archive, weeks, today])

  if (!data) return <p className="p-6 text-sm text-muted-foreground">Loading…</p>

  const filterBoard = (search: Record<string, string>) => void toBoard({ to: '/', search })
  const dot = (color: string) => (
    <span className="size-2 shrink-0 rounded-full" style={{ backgroundColor: color }} />
  )

  return (
    <main className="flex-1 overflow-y-auto">
      <div className="mx-auto max-w-6xl space-y-4 p-4 sm:p-6">
        <div className="flex flex-wrap items-center gap-3">
          <h1 className="text-lg font-semibold tracking-tight">Insights</h1>
          <Tabs
            value={range}
            onValueChange={(v) =>
              void navigate({ search: { range: v as RangeKey }, replace: true })
            }
            className="ml-auto"
          >
            <TabsList>
              {Object.keys(ranges).map((key) => (
                <TabsTrigger key={key} value={key} className="px-3 text-xs">
                  {key}
                </TabsTrigger>
              ))}
            </TabsList>
          </Tabs>
        </div>

        <div className="grid grid-cols-2 gap-3 md:grid-cols-3 xl:grid-cols-6">
          <StatTile label="Created" value={view.created} hint={`last ${weeks} weeks`} />
          <StatTile label="Claimed" value={view.claimed} hint={`last ${weeks} weeks`} />
          <StatTile label="Archived as done" value={view.done} hint={`last ${weeks} weeks`} />
          <StatTile label="Open now" value={view.open} hint="before merged" />
          <StatTile label="In progress now" value={view.inProgress} />
          <StatTile label="Blocked now" value={view.blocked} />
        </div>

        <WeeklyChart weeks={view.weekly} />

        <div className="grid gap-4 lg:grid-cols-2">
          <RankedBars
            title="Tasks by stage"
            description="Where every task on the board sits now."
            slices={view.byStage}
            labelOf={(key) => (
              <>
                <StatusIcon
                  status={view.byStage.find((s) => s.key === key)?.status ?? ''}
                  className="size-3"
                />
                {key}
              </>
            )}
          />
          <RankedBars
            title="Open work by epic"
            description="Tasks not merged yet, per epic."
            slices={view.byEpic}
            labelOf={(key) => (
              <>
                {dot(colorOf(`epic:${key}`))}
                {key}
              </>
            )}
            onSelect={(key) => filterBoard({ label: `epic:${key}` })}
          />
          <RankedBars
            title="Open work by type"
            slices={view.byType}
            labelOf={(key) => (
              <>
                {dot(colorOf(`type:${key}`))}
                {key}
              </>
            )}
            onSelect={(key) => filterBoard({ label: `type:${key}` })}
          />
          <RankedBars
            title="Claimed open work by owner"
            slices={view.byOwner}
            labelOf={(key) => (
              <>
                <Avatar name={key} className="size-4 text-[7px]" />
                {key}
              </>
            )}
            onSelect={(key) => filterBoard({ owner: key })}
          />
        </div>
      </div>
    </main>
  )
}
