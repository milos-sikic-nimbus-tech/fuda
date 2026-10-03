import { createFileRoute, Link } from '@tanstack/react-router'
import { RichContent } from '@/components/molecules/RichContent'
import { useBoard, useGuide, useGuidePage } from '@/lib/queries'
import { cn } from '@/lib/utils'

export const Route = createFileRoute('/guide/$')({
  component: GuidePage,
})

function GuidePage() {
  const { _splat: splat = '' } = Route.useParams()
  const { data: pages = [] } = useGuide()
  const slug = splat || pages[0]?.slug || 'overview'
  const { data: page, error } = useGuidePage(slug)
  const { data: board } = useBoard()

  return (
    <main className="flex min-h-0 flex-1 flex-col md:flex-row">
      <nav className="flex shrink-0 gap-1 overflow-x-auto border-b border-border/60 px-3 py-2 md:w-56 md:flex-col md:overflow-visible md:border-r md:border-b-0 md:px-3 md:py-6">
        <p className="hidden px-2 pb-2 text-xs font-medium text-muted-foreground md:block">Guide</p>
        {pages.map((p) => (
          <Link
            key={p.slug}
            to="/guide/$"
            params={{ _splat: p.slug }}
            className={cn(
              'shrink-0 rounded-md px-2 py-1.5 text-sm whitespace-nowrap text-muted-foreground transition-colors hover:bg-accent hover:text-foreground',
              p.slug === slug && 'bg-accent font-medium text-foreground',
            )}
          >
            {p.title}
          </Link>
        ))}
      </nav>
      <div className="min-w-0 flex-1 overflow-y-auto">
        <article className="mx-auto max-w-3xl px-4 py-6 sm:px-8">
          {error && (
            <p className="text-sm text-muted-foreground">This guide page does not exist.</p>
          )}
          {page && <RichContent html={page.html} taskIds={board?.cards.map((c) => c.id) ?? []} />}
        </article>
      </div>
    </main>
  )
}
