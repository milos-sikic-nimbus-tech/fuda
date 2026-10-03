import { createFileRoute, Link } from '@tanstack/react-router'
import { useMemo, useState } from 'react'
import { Input } from '@/components/ui/input'
import { useArchive } from '@/lib/queries'

export const Route = createFileRoute('/archive')({
  component: ArchivePage,
})

function ArchivePage() {
  const { data = [], isLoading } = useArchive()
  const [q, setQ] = useState('')
  const shown = useMemo(() => {
    const needle = q.trim().toLowerCase()
    return needle ? data.filter((c) => `${c.id} ${c.title}`.toLowerCase().includes(needle)) : data
  }, [data, q])

  return (
    <main className="mx-auto w-full max-w-4xl flex-1 overflow-y-auto p-6">
      <div className="mb-4 flex items-center gap-3">
        <h1 className="text-lg font-semibold">Archive</h1>
        <Input
          value={q}
          onChange={(e) => setQ(e.target.value)}
          placeholder="Search archived tasks"
          className="h-8 w-64"
        />
        <span className="ml-auto text-xs text-muted-foreground">{shown.length} tasks</span>
      </div>
      {isLoading ? null : shown.length === 0 ? (
        <p className="text-sm text-muted-foreground">
          Nothing archived{q ? ' matches' : ' yet'}. Tasks move to docs/board/archive/ after they
          have been in prod for a while.
        </p>
      ) : (
        <ul className="divide-y divide-border rounded-lg border border-border bg-card">
          {shown.map((c) => (
            <li key={c.id}>
              <Link
                to="/"
                search={{ task: c.id }}
                className="flex items-center gap-3 px-4 py-2.5 text-sm hover:bg-accent"
              >
                <span className="w-20 font-mono text-xs text-muted-foreground">{c.id}</span>
                <span className="flex-1">{c.title}</span>
                <span className="text-xs text-muted-foreground">{c.done}</span>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </main>
  )
}
