import type { ReactNode } from 'react'
import { Ban, GitPullRequest } from 'lucide-react'
import { LabelChip } from '@/components/atoms/LabelChip'
import { PersonChip } from '@/components/atoms/PersonChip'
import { StatusIcon } from '@/components/atoms/StatusIcon'
import { TaskLink } from '@/components/atoms/TaskLink'
import { RichContent } from '@/components/molecules/RichContent'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { useBoardSearch } from '@/hooks/useBoardSearch'
import { listOf } from '@/lib/filters'
import { useTask } from '@/lib/queries'
import { prHref } from '@/lib/api'

function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <>
      <dt className="pt-0.5 text-xs text-muted-foreground">{label}</dt>
      <dd className="flex min-w-0 flex-wrap items-center gap-1.5 text-[13px] [overflow-wrap:anywhere]">
        {children}
      </dd>
    </>
  )
}

function blockedNote(text: string, ids: string[]): boolean {
  return ids.reduce((rest, id) => rest.replaceAll(id, ''), text).replace(/[\s,;]/g, '') !== ''
}

function Ids({ ids }: { ids: string[] }) {
  return ids.map((id) => <TaskLink key={id} id={id} />)
}

export function TaskSheet({ taskIds, prLink }: { taskIds: string[]; prLink?: string }) {
  const { search, toggle, openTask } = useBoardSearch()
  const { data: task, isError } = useTask(search.task)
  const labels = listOf(search.label)

  return (
    <Sheet open={!!search.task} onOpenChange={(open) => !open && openTask(undefined)}>
      <SheetContent className="w-full gap-0 overflow-y-auto data-[side=right]:w-full data-[side=right]:sm:max-w-3xl">
        <SheetHeader className="border-b">
          <SheetDescription className="flex items-center gap-2 font-mono text-xs">
            {search.task}
            {task && (
              <span className="inline-flex items-center gap-1 font-sans">
                <StatusIcon status={task.status} className="size-3" />
                {task.status}
              </span>
            )}
          </SheetDescription>
          <SheetTitle className="text-lg leading-snug">
            {task?.title ?? (isError ? 'Task not found' : 'Loading…')}
          </SheetTitle>
        </SheetHeader>
        {task && (
          <div className="space-y-5 p-4">
            <dl className="grid grid-cols-[7rem_1fr] gap-x-3 gap-y-2">
              {task.owners.length > 0 && (
                <Field label="Owner">
                  {task.owners.map((o) => (
                    <PersonChip key={o} name={o} onClick={() => toggle('owner', o)} />
                  ))}
                </Field>
              )}
              {task.testers.length > 0 && (
                <Field label="Tester">
                  {task.testers.map((o) => (
                    <PersonChip key={o} name={o} onClick={() => toggle('tester', o)} />
                  ))}
                </Field>
              )}
              {task.labels.length > 0 && (
                <Field label="Labels">
                  {task.labels.map((l) => (
                    <LabelChip
                      key={l}
                      label={l}
                      active={labels.includes(l)}
                      onClick={() => toggle('label', l)}
                    />
                  ))}
                </Field>
              )}
              {task.blockedBy && (
                <Field label="Blocked by">
                  <Ban className="size-3.5 text-amber-700" />
                  <Ids ids={task.blockedByIds} />
                  {blockedNote(task.blockedBy, task.blockedByIds) && (
                    <span className="text-xs text-muted-foreground">{task.blockedBy}</span>
                  )}
                </Field>
              )}
              {task.blocks.length > 0 && (
                <Field label="Blocks">
                  <Ids ids={task.blocks} />
                </Field>
              )}
              {(task.prs.length > 0 || task.prRef || task.openPrs.length > 0) && (
                <Field label="PR">
                  {[...new Set([...task.prs, ...task.openPrs])].map((n) => (
                    <a
                      key={n}
                      href={prHref(prLink, n)}
                      target="_blank"
                      rel="noreferrer"
                      className="inline-flex items-center gap-1 text-xs text-primary hover:underline"
                    >
                      <GitPullRequest className="size-3.5" />#{n}
                      {task.openPrs.includes(n) && <span className="text-violet-700">open</span>}
                    </a>
                  ))}
                  {task.prRef && <span className="font-mono text-xs">{task.prRef}</span>}
                </Field>
              )}
              {task.added && <Field label="Added">{task.added}</Field>}
              {task.claimed && <Field label="Claimed">{task.claimed}</Field>}
              {task.done && <Field label="Done">{task.done}</Field>}
              {task.inProd && (
                <Field label="Prod">in prod{task.archiveCandidate && ' · ready to archive'}</Field>
              )}
              {task.custom.map((f) => (
                <Field key={f.key} label={f.key}>
                  <span className="text-muted-foreground">{f.value || '—'}</span>
                </Field>
              ))}
              {task.referencedBy.length > 0 && (
                <Field label="Referenced by">
                  <Ids ids={task.referencedBy} />
                </Field>
              )}
            </dl>
            <RichContent
              html={task.html}
              taskIds={taskIds}
              selfId={task.id}
              lead="Done when"
              dropTitle
            />
            <p className="font-mono text-[11px] text-muted-foreground">{task.path}</p>
          </div>
        )}
      </SheetContent>
    </Sheet>
  )
}
