import { useBoard } from '@/lib/queries'

const fallback = '#94a3b8'

export function useLabelColor(): (label: string) => string {
  const { data } = useBoard()
  const colors = data?.facets.labelColors
  return (label) => colors?.[label] ?? fallback
}
