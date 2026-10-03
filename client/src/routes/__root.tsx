import type { QueryClient } from '@tanstack/react-query'
import { createRootRouteWithContext, Link, Outlet } from '@tanstack/react-router'
import { useEffect } from 'react'
import { SyncIndicator } from '@/components/molecules/SyncIndicator'
import { ProblemsDrawer } from '@/components/organisms/ProblemsDrawer'
import { TooltipProvider } from '@/components/ui/tooltip'
import { useBoard } from '@/lib/queries'

export const Route = createRootRouteWithContext<{ queryClient: QueryClient }>()({
  component: RootLayout,
})

const navClass =
  'rounded-md px-2.5 py-1 text-sm text-muted-foreground hover:text-foreground data-[status=active]:bg-accent data-[status=active]:text-foreground'

function RootLayout() {
  const { data } = useBoard()

  useEffect(() => {
    document.title = data?.title ? `${data.title} · fuda` : 'fuda'
  }, [data?.title])

  return (
    <TooltipProvider delayDuration={300}>
      <div className="flex h-svh flex-col text-foreground">
        <header className="flex items-center gap-3 border-b border-border/60 bg-card/70 px-4 py-2 backdrop-blur-md">
          <Link to="/" className="flex items-baseline gap-2 font-semibold">
            <span className="text-primary">fuda</span>
            {data?.title && (
              <span className="text-sm font-medium text-muted-foreground">{data.title}</span>
            )}
          </Link>
          <nav className="ml-4 flex items-center gap-1">
            <Link to="/" className={navClass} activeOptions={{ exact: true, includeSearch: false }}>
              Board
            </Link>
            <Link to="/archive" className={navClass}>
              Archive
            </Link>
            <Link to="/docs/$" params={{ _splat: 'board/README.md' }} className={navClass}>
              Docs
            </Link>
            <Link to="/guide/$" params={{ _splat: '' }} className={navClass}>
              Guide
            </Link>
          </nav>
          <div className="ml-auto flex items-center gap-3">
            {data && <ProblemsDrawer problems={data.problems} />}
            {data && <SyncIndicator sync={data.sync} />}
          </div>
        </header>
        <Outlet />
      </div>
    </TooltipProvider>
  )
}
