import { createFileRoute } from '@tanstack/react-router'
import { RichContent } from '@/components/molecules/RichContent'
import { TaskLink } from '@/components/atoms/TaskLink'
import { DocTree } from '@/components/organisms/DocTree'
import { useBoard, useDoc, useDocList } from '@/lib/queries'

export const Route = createFileRoute('/docs/$')({
  component: DocPage,
})

const defaultDoc = 'board/README.md'

function DocPage() {
  const { _splat: requested = '' } = Route.useParams()
  const { data: paths } = useDocList()
  const { data: board } = useBoard()
  const path = requested || (paths?.includes(defaultDoc) ? defaultDoc : '')
  const { data: doc, error } = useDoc(path)

  return (
    <div className="flex min-h-0 flex-1">
      <DocTree paths={paths ?? []} openPath={path} />
      <main className="flex-1 overflow-y-auto">
        <div className="mx-auto max-w-3xl p-6">
          {path === '' && paths && <p className="text-sm text-muted-foreground">Pick a doc</p>}
          {path !== '' && (
            <p className="mb-4 font-mono text-xs text-muted-foreground">docs/{path}</p>
          )}
          {error && (
            <p className="text-sm text-muted-foreground">This doc does not exist in the repo.</p>
          )}
          {doc && (
            <>
              <RichContent html={doc.html} taskIds={board?.cards.map((c) => c.id) ?? []} />
              {doc.backlinks.length > 0 && (
                <footer className="mt-8 flex flex-wrap items-center gap-2 border-t border-border pt-4 text-xs text-muted-foreground">
                  Mentioned by{' '}
                  {doc.backlinks.map((id) => (
                    <TaskLink key={id} id={id} />
                  ))}
                </footer>
              )}
            </>
          )}
        </div>
      </main>
    </div>
  )
}
