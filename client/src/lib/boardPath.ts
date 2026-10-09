const segmentCounts: Record<string, number> = { github: 2, azure: 3, local: 1 }

export function boardPrefix(pathname: string): string {
  const [host, ...rest] = pathname.split('/').filter(Boolean)
  const count = segmentCounts[host ?? '']
  if (!count || rest.length < count) return ''
  return '/' + [host, ...rest.slice(0, count)].join('/')
}

export const BOARD_PREFIX =
  typeof window === 'undefined' ? '' : boardPrefix(window.location.pathname)
