import { Archive, Ban, GitPullRequest, Rocket, Sparkles } from 'lucide-react'
import type { ReactNode } from 'react'
import { Hint } from '@/components/atoms/Hint'
import { LabelChip, labelText } from '@/components/atoms/LabelChip'
import { Avatar } from '@/components/atoms/PersonChip'
import { StatusDot } from '@/components/atoms/StatusDot'
import { useBoardSearch } from '@/hooks/useBoardSearch'
import { type Card, type Column, prHref } from '@/lib/api'
import { hueStyle, statusHue } from '@/lib/colors'
import { listOf } from '@/lib/filters'
import { cn } from '@/lib/utils'

const visibleLabels = 2
const labelPriority = ['type', 'epic']

function rankLabels(labels: string[]): string[] {
  const rank = (label: string) => {
    const i = labelPriority.indexOf(labelText(label).group)
    return i === -1 ? labelPriority.length : i
  }
  return [...labels].sort((a, b) => rank(a) - rank(b))
}

function Badge({
  hint,
  className,
  children,
}: {
  hint: ReactNode
  className: string
  children: ReactNode
}) {
  return (
    <Hint label={hint}>
      <span
        className={cn(
          'inline-flex max-w-full items-center gap-1 rounded-md px-1.5 py-0.5 text-[11px] leading-tight',
          className,
        )}
      >
        {children}
      </span>
    </Hint>
  )
}

function CardBadges({ card, prLink }: { card: Card; prLink?: string }) {
  return (
    <>
      {card.blockedBy && (
        <Badge
          hint={`Blocked by ${card.blockedBy}`}
          className="max-w-56 bg-amber-100 text-amber-900"
        >
          <Ban className="size-3 shrink-0" />
          <span className="truncate">{card.blockedBy}</span>
        </Badge>
      )}
      {card.newInPr && (
        <Badge
          hint="The task file is new in an open PR and not on develop yet"
          className="bg-violet-50 text-violet-700"
        >
          <Sparkles className="size-3" />
          new in PR
        </Badge>
      )}
      {card.openPrs.map((n) => (
        <Hint key={n} label={`Open pull request #${n} delivers this task`}>
          <a
            href={prHref(prLink, n)}
            target="_blank"
            rel="noreferrer"
            onClick={(e) => e.stopPropagation()}
            className="inline-flex items-center gap-1 rounded-md bg-violet-100 px-1.5 py-0.5 text-[11px] leading-tight text-violet-800 hover:bg-violet-200"
          >
            <GitPullRequest className="size-3" />#{n}
          </a>
        </Hint>
      ))}
      {card.inProd && (
        <Badge
          hint="Its file on main is merged, testing or validated"
          className="bg-emerald-100 text-emerald-800"
        >
          <Rocket className="size-3" />
          in prod
        </Badge>
      )}
      {card.archiveCandidate && (
        <Badge
          hint="In prod long enough to move to the archive"
          className="bg-zinc-100 text-zinc-700"
        >
          <Archive className="size-3" />
          archive
        </Badge>
      )}
    </>
  )
}

export function TaskCard({
  card,
  column,
  prLink,
}: {
  card: Card
  column: Column
  prLink?: string
}) {
  const { search, toggle, openTask } = useBoardSearch()
  const selectedLabels = listOf(search.label)
  const ranked = rankLabels(card.labels)
  const hiddenLabels = ranked.slice(visibleLabels)
  const showStatus = column.statuses.length > 1 || column.prOpen

  return (
    <div
      role="button"
      tabIndex={0}
      onClick={() => openTask(card.id)}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault()
          openTask(card.id)
        }
      }}
      style={hueStyle(statusHue(column.prOpen ? 'in review' : card.status))}
      className={cn(
        'cursor-pointer rounded-lg border border-border/80 bg-card p-3 text-left shadow-xs transition duration-150 hue-lift',
        'focus-visible:ring-2 focus-visible:ring-ring/60 focus-visible:outline-none',
        search.task === card.id && 'border-primary ring-2 ring-primary/25',
      )}
    >
      <div className="mb-1 flex items-center gap-2 text-[11px] text-muted-foreground">
        <span className="font-mono font-medium tracking-tight text-foreground/70">{card.id}</span>
        {showStatus && (
          <Hint label="This column holds several statuses; this is the task's own">
            <span className="inline-flex items-center gap-1">
              <StatusDot status={card.status} className="size-1.5" />
              {card.status}
            </span>
          </Hint>
        )}
        <span className="ml-auto flex -space-x-1.5">
          {card.owners.map((owner) => (
            <Avatar
              key={owner}
              name={owner}
              hint={`Owner: ${owner}`}
              className="ring-2 ring-card"
            />
          ))}
        </span>
      </div>
      <p className="text-[13px] leading-snug font-medium text-card-foreground">{card.title}</p>
      <div className="mt-2 flex flex-wrap items-center gap-1">
        <CardBadges card={card} prLink={prLink} />
        {ranked.slice(0, visibleLabels).map((label) => (
          <LabelChip
            key={label}
            label={label}
            showGroup={false}
            active={selectedLabels.includes(label)}
            onClick={() => toggle('label', label)}
          />
        ))}
        {hiddenLabels.length > 0 && (
          <Hint label={hiddenLabels.join(', ')}>
            <span className="text-[11px] text-muted-foreground">+{hiddenLabels.length}</span>
          </Hint>
        )}
      </div>
    </div>
  )
}
