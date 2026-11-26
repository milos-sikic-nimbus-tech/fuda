import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { cleanup, render } from '@testing-library/react'
import type { ReactElement } from 'react'
import { afterEach } from 'vitest'
import { keys } from '@/lib/queries'
import type { BoardData } from '@/lib/api'
import type { BoardListData } from '@/lib/hosts'

afterEach(cleanup)

export function renderWithListing(ui: ReactElement, listing: BoardListData) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  client.setQueryData(keys.boards, listing)

  return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>)
}

export function renderWithBoard(ui: ReactElement, board: Partial<BoardData> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  client.setQueryData(keys.board, { facets: { people: [] }, ...board })

  return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>)
}
