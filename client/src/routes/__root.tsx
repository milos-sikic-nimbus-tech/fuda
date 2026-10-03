import type { QueryClient } from '@tanstack/react-query'
import { createRootRouteWithContext, Outlet } from '@tanstack/react-router'
import { useEffect } from 'react'
import { AppHeader } from '@/components/organisms/AppHeader'
import { Toaster } from '@/components/ui/sonner'
import { TooltipProvider } from '@/components/ui/tooltip'
import { useTheme } from '@/hooks/useTheme'
import { useBoard } from '@/lib/queries'

export const Route = createRootRouteWithContext<{ queryClient: QueryClient }>()({
  component: RootLayout,
})

function RootLayout() {
  const { data } = useBoard()
  const { resolved } = useTheme()

  useEffect(() => {
    document.title = data?.title ? `${data.title} · fuda` : 'fuda'
  }, [data?.title])

  return (
    <TooltipProvider delayDuration={300}>
      <div className="flex h-svh flex-col text-foreground">
        <AppHeader board={data} />
        <Outlet />
      </div>
      <Toaster theme={resolved} position="bottom-right" richColors closeButton />
    </TooltipProvider>
  )
}
