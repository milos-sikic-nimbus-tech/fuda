import { keepPreviousData, useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { ApiError, api } from './api'
import { BOARD_PREFIX } from './boardPath'

export const keys = {
  boards: ['boards'] as const,
  board: ['board'] as const,
  task: (id: string) => ['task', id] as const,
  search: (q: string) => ['search', q] as const,
  archive: ['archive'] as const,
  doc: (path: string) => ['doc', path] as const,
  guide: ['guide'] as const,
  guidePage: (slug: string) => ['guide', slug] as const,
}

export function useBoards() {
  return useQuery({ queryKey: keys.boards, queryFn: api.boards, staleTime: 5 * 60_000 })
}

export function useBoard() {
  return useQuery({
    queryKey: keys.board,
    queryFn: api.board,
    refetchInterval: 5_000,
    staleTime: 0,
    enabled: !!BOARD_PREFIX,
  })
}

export function useTask(id: string | undefined) {
  return useQuery({
    queryKey: keys.task(id ?? ''),
    queryFn: () => api.task(id!),
    enabled: !!id,
  })
}

export function useTextSearch(q: string | undefined, enabled: boolean) {
  const query = q?.trim() ?? ''
  return useQuery({
    queryKey: keys.search(query),
    queryFn: () => api.search(query),
    enabled: enabled && query.length > 1,
    placeholderData: keepPreviousData,
  })
}

export function useArchive() {
  return useQuery({ queryKey: keys.archive, queryFn: api.archive })
}

export function useDoc(path: string) {
  return useQuery({ queryKey: keys.doc(path), queryFn: () => api.doc(path) })
}

export function useSync() {
  const client = useQueryClient()
  return useMutation({
    mutationFn: api.sync,
    onSuccess: () => {
      toast.success('Reading the repo…')
      setTimeout(() => void client.invalidateQueries(), 1500)
    },
    onError: (error) => {
      if (error instanceof ApiError && error.status === 429) {
        toast.info('Synced moments ago', { description: 'Try again in a few seconds.' })
        return
      }
      toast.error('Sync failed', { description: error.message })
    },
  })
}

export function useGuide() {
  return useQuery({ queryKey: keys.guide, queryFn: api.guide, staleTime: Infinity })
}

export function useGuidePage(slug: string) {
  return useQuery({
    queryKey: keys.guidePage(slug),
    queryFn: () => api.guidePage(slug),
    staleTime: Infinity,
  })
}
