import { Link } from '@tanstack/react-router'
import { Archive, BookOpen, ChartColumn, FileText, LayoutGrid, type LucideIcon } from 'lucide-react'
import type { ReactNode } from 'react'
import { HostIcon } from '@/components/atoms/HostIcon'
import { Logo } from '@/components/atoms/Logo'
import { SyncIndicator } from '@/components/molecules/SyncIndicator'
import { ThemeToggle } from '@/components/molecules/ThemeToggle'
import { ProblemsDrawer } from '@/components/organisms/ProblemsDrawer'
import type { BoardData } from '@/lib/api'
import { cn } from '@/lib/utils'

const tabClass = cn(
  'relative flex h-full items-center gap-1.5 px-2 text-sm text-muted-foreground transition-colors hover:text-foreground sm:px-2.5',
  'data-[status=active]:text-foreground data-[status=active]:after:absolute data-[status=active]:after:inset-x-2',
  'data-[status=active]:after:-bottom-px data-[status=active]:after:h-0.5 data-[status=active]:after:rounded-full data-[status=active]:after:bg-primary',
)

function TabLabel({ icon: Icon, children }: { icon: LucideIcon; children: ReactNode }) {
  return (
    <>
      <Icon className="size-4 sm:hidden" aria-hidden="true" />
      <span className="sr-only sm:not-sr-only">{children}</span>
    </>
  )
}

function OriginPill({ board }: { board: BoardData }) {
  const content = (
    <>
      <HostIcon host={board.origin.host} className="size-3.5 shrink-0" />
      <span className="truncate font-medium text-foreground">
        {board.title || board.origin.repo}
      </span>
      <span className="hidden text-muted-foreground md:inline">· {board.sync.develop.branch}</span>
    </>
  )
  const className =
    'inline-flex min-w-0 max-w-40 items-center gap-1.5 rounded-full border border-border bg-card px-2.5 py-0.5 text-[13px] shadow-xs md:max-w-none'
  return board.origin.url ? (
    <a
      href={board.origin.url}
      target="_blank"
      rel="noreferrer"
      className={cn(className, 'hover:border-foreground/20')}
    >
      {content}
    </a>
  ) : (
    <span className={className}>{content}</span>
  )
}

export function AppHeader({ board }: { board?: BoardData }) {
  return (
    <header className="flex h-12 shrink-0 items-center gap-2 border-b border-border/70 bg-card/80 px-3 backdrop-blur-md sm:gap-3 sm:px-4">
      <Link to="/" className="flex shrink-0 items-center gap-2" aria-label="fuda">
        <Logo className="size-6" />
        <span className="hidden text-[15px] font-semibold tracking-tight sm:inline">fuda</span>
      </Link>
      {board && (
        <>
          <span className="hidden text-border sm:inline">/</span>
          <OriginPill board={board} />
        </>
      )}
      <nav className="flex h-full min-w-0 items-stretch sm:ml-2">
        <Link to="/" className={tabClass} activeOptions={{ exact: true, includeSearch: false }}>
          <TabLabel icon={LayoutGrid}>Board</TabLabel>
        </Link>
        <Link to="/insights" className={tabClass}>
          <TabLabel icon={ChartColumn}>Insights</TabLabel>
        </Link>
        <Link to="/archive" className={tabClass}>
          <TabLabel icon={Archive}>Archive</TabLabel>
        </Link>
        <Link to="/docs/$" params={{ _splat: 'board/README.md' }} className={tabClass}>
          <TabLabel icon={FileText}>Docs</TabLabel>
        </Link>
        <Link to="/guide/$" params={{ _splat: '' }} className={tabClass}>
          <TabLabel icon={BookOpen}>Guide</TabLabel>
        </Link>
      </nav>
      <div className="ml-auto flex shrink-0 items-center gap-2 sm:gap-3">
        {board && <ProblemsDrawer problems={board.problems} />}
        {board && <SyncIndicator sync={board.sync} />}
        <ThemeToggle />
      </div>
    </header>
  )
}
