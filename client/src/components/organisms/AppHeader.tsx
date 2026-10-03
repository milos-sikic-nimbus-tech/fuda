import { Link } from '@tanstack/react-router'
import { HostIcon } from '@/components/atoms/HostIcon'
import { Logo } from '@/components/atoms/Logo'
import { SyncIndicator } from '@/components/molecules/SyncIndicator'
import { ThemeToggle } from '@/components/molecules/ThemeToggle'
import { ComingSoon } from '@/components/organisms/ComingSoon'
import { ProblemsDrawer } from '@/components/organisms/ProblemsDrawer'
import type { BoardData } from '@/lib/api'
import { cn } from '@/lib/utils'

const tabClass =
  'relative flex h-full items-center px-2.5 text-sm text-muted-foreground transition-colors hover:text-foreground data-[status=active]:text-foreground data-[status=active]:after:absolute data-[status=active]:after:inset-x-2 data-[status=active]:after:-bottom-px data-[status=active]:after:h-0.5 data-[status=active]:after:rounded-full data-[status=active]:after:bg-primary'

function OriginPill({ board }: { board: BoardData }) {
  const content = (
    <>
      <HostIcon host={board.origin.host} className="size-3.5" />
      <span className="font-medium text-foreground">{board.title || board.origin.repo}</span>
      <span className="text-muted-foreground">· {board.sync.develop.branch}</span>
    </>
  )
  const className =
    'inline-flex items-center gap-1.5 rounded-full border border-border bg-card px-2.5 py-0.5 text-[13px] shadow-xs'
  return board.origin.url ? (
    <a
      href={board.origin.url}
      target="_blank"
      rel="noreferrer"
      className={cn(className, 'hover:border-primary/40')}
    >
      {content}
    </a>
  ) : (
    <span className={className}>{content}</span>
  )
}

export function AppHeader({ board }: { board?: BoardData }) {
  return (
    <header className="flex h-12 shrink-0 items-center gap-3 border-b border-border/70 bg-card/80 px-4 backdrop-blur-md">
      <Link to="/" className="flex items-center gap-2">
        <Logo className="size-6" />
        <span className="text-[15px] font-semibold tracking-tight">fuda</span>
      </Link>
      {board && (
        <>
          <span className="text-border">/</span>
          <OriginPill board={board} />
        </>
      )}
      <nav className="ml-4 flex h-full items-stretch">
        <Link to="/" className={tabClass} activeOptions={{ exact: true, includeSearch: false }}>
          Board
        </Link>
        <Link to="/archive" className={tabClass}>
          Archive
        </Link>
        <Link to="/docs/$" params={{ _splat: 'board/README.md' }} className={tabClass}>
          Docs
        </Link>
        <ComingSoon
          feature="The guide"
          description="How to set up a repo for fuda, the task format, board config files and the agent guide will live here. Until then, see the repo's docs/board/TASKS.md."
        >
          <button type="button" className={tabClass}>
            Guide
          </button>
        </ComingSoon>
      </nav>
      {board && (
        <div className="ml-auto flex items-center gap-3">
          <ProblemsDrawer problems={board.problems} />
          <SyncIndicator sync={board.sync} />
        </div>
      )}
      <div className={board ? '' : 'ml-auto'}>
        <ThemeToggle />
      </div>
    </header>
  )
}
