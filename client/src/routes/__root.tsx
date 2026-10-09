import type { QueryClient } from '@tanstack/react-query'
import { createRootRouteWithContext, Outlet } from '@tanstack/react-router'
import { useEffect } from 'react'
import { AppHeader } from '@/components/organisms/AppHeader'
import { Toaster } from '@/components/ui/sonner'
import { TooltipProvider } from '@/components/ui/tooltip'
import { useTheme } from '@/hooks/useTheme'
import { ApiError } from '@/lib/api'
import { useBoard } from '@/lib/queries'

export const Route = createRootRouteWithContext<{ queryClient: QueryClient }>()({
  component: RootLayout,
})

function RootLayout() {
  const { data, error } = useBoard()
  const { resolved } = useTheme()

  useEffect(() => {
    document.title = data?.title ? `${data.title} · fuda` : 'fuda'
  }, [data?.title])

  return (
    <TooltipProvider delayDuration={300}>
      <div className="flex h-svh flex-col text-foreground">
        <AppHeader board={data} />
        {error instanceof ApiError && error.status === 404 ? (
          <p className="p-6 text-sm text-muted-foreground">No Board at this address.</p>
        ) : (
          <Outlet />
        )}
      </div>
      <Toaster theme={resolved} position="bottom-right" richColors closeButton />
    </TooltipProvider>
  )
}
