import { QueryCache, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createRouter, RouterProvider } from '@tanstack/react-router'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { toast } from 'sonner'
import { ApiError } from './lib/api'
import { BOARD_PREFIX } from './lib/boardPath'
import { routeTree } from './routeTree.gen'
import './index.css'

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 30_000, refetchOnWindowFocus: true } },
  queryCache: new QueryCache({
    onError: (error) => {
      if (error instanceof ApiError && error.status === 404) return
      toast.error('Could not load data', { description: error.message })
    },
  }),
})

const router = createRouter({ routeTree, basepath: BOARD_PREFIX || '/', context: { queryClient } })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
)
