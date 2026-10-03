import { Ban, GitPullRequest, Rocket, Archive } from 'lucide-react'
import { LabelChip } from '@/components/atoms/LabelChip'
import { Avatar } from '@/components/atoms/PersonChip'
import { StatusDot } from '@/components/atoms/StatusDot'
import { useBoardSearch } from '@/hooks/useBoardSearch'
import type { Card, Column } from '@/lib/api'
import { listOf } from '@/lib/filters'
import { cn } from '@/lib/utils'

const maxLabels = 3

export function TaskCard({ card, column }: { card: Card; column: Column }) {
  const { search, toggle, openTask } = useBoardSearch()
  const selectedLabels = listOf(search.label)
  const showStatus = column.statuses.length > 1 || column.prOpen
  const extraLabels = card.labels.length - maxLabels

  return (
    <article
      onClick={() => openTask(card.id)}
      className={cn(
        'group cursor-pointer rounded-lg border border-border bg-card p-3 shadow-xs transition hover:-translate-y-px hover:border-primary/40 hover:shadow-sm',
        search.task === card.id && 'border-primary ring-2 ring-primary/20',
      )}
    >
      <div className="mb-1.5 flex items-center gap-2 text-xs text-muted-foreground">
        <span className="font-mono font-medium text-foreground/80">{card.id}</span>
        {showStatus && (
          <span className="inline-flex items-center gap-1">
            <StatusDot status={card.status} />
            {card.status}
          </span>
        )}
        <span className="ml-auto flex -space-x-1">
          {card.owners.map((o) => (
            <Avatar key={o} name={o} className="ring-2 ring-card" />
          ))}
        </span>
      </div>
      <p className="text-sm leading-snug text-card-foreground">{card.title}</p>
      {(card.labels.length > 0 || card.blockedBy || card.openPrs.length > 0 || card.inProd) && (
        <div className="mt-2 flex flex-wrap items-center gap-1">
          {card.blockedBy && (
            <span
              title={`Blocked by ${card.blockedBy}`}
              className="inline-flex max-w-full items-center gap-1 truncate rounded-md bg-amber-100 px-1.5 py-0.5 text-[11px] text-amber-800"
            >
              <Ban className="size-3 shrink-0" />
              <span className="truncate">{card.blockedBy}</span>
            </span>
          )}
          {card.openPrs.map((n) => (
            <span
              key={n}
              className="inline-flex items-center gap-1 rounded-md bg-violet-100 px-1.5 py-0.5 text-[11px] text-violet-800"
            >
              <GitPullRequest className="size-3" />#{n}
            </span>
          ))}
          {card.inProd && (
            <span className="inline-flex items-center gap-1 rounded-md bg-emerald-100 px-1.5 py-0.5 text-[11px] text-emerald-800">
              <Rocket className="size-3" />
              in prod
            </span>
          )}
          {card.archiveCandidate && (
            <span className="inline-flex items-center gap-1 rounded-md bg-zinc-100 px-1.5 py-0.5 text-[11px] text-zinc-700">
              <Archive className="size-3" />
              archive
            </span>
          )}
          {card.labels.slice(0, maxLabels).map((l) => (
            <LabelChip
              key={l}
              label={l}
              active={selectedLabels.includes(l)}
              onClick={() => toggle('label', l)}
            />
          ))}
          {extraLabels > 0 && (
            <span className="text-[11px] text-muted-foreground">+{extraLabels}</span>
          )}
        </div>
      )}
    </article>
  )
}
