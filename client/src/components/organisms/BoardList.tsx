import { useEffect } from 'react'
import { HostIcon } from '@/components/atoms/HostIcon'
import { useBoards } from '@/lib/queries'

export function BoardList() {
  const { data: boards, error } = useBoards()
  const only = boards?.length === 1 ? boards[0] : undefined

  useEffect(() => {
    if (only) window.location.replace(only.path)
  }, [only])

  if (error) {
    return (
      <p className="p-6 text-sm text-muted-foreground">Could not list Boards: {error.message}</p>
    )
  }
  if (!boards || only) return <p className="p-6 text-sm text-muted-foreground">Loading Boards…</p>
  if (boards.length === 0) {
    return (
      <p className="p-6 text-sm text-muted-foreground">
        No Boards yet. Create a repository whose name starts with <code>fuda-</code> (on GitHub,
        install the fuda GitHub App on it).
      </p>
    )
  }
  return (
    <main className="mx-auto w-full max-w-md p-6">
      <h1 className="mb-3 text-base font-semibold">Boards</h1>
      <ul className="flex flex-col gap-1">
        {boards.map((board) => (
          <li key={board.path}>
            <a
              href={board.path}
              className="flex items-center gap-2 rounded-md border border-border bg-card px-3 py-2 text-sm hover:border-foreground/20"
            >
              <HostIcon host={board.host} className="size-4 shrink-0" />
              <span className="truncate">{board.repo}</span>
            </a>
          </li>
        ))}
      </ul>
    </main>
  )
}
