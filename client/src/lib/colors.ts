const statusColors: Record<string, string> = {
  backlog: 'bg-slate-400',
  'in progress': 'bg-blue-500',
  'in review': 'bg-violet-500',
  merged: 'bg-teal-500',
  testing: 'bg-amber-500',
  validated: 'bg-emerald-600',
  done: 'bg-emerald-800',
}

export function statusColor(status: string): string {
  return statusColors[status.trim().toLowerCase()] ?? 'bg-zinc-300'
}

const hues = [
  '#10b981',
  '#6366f1',
  '#f59e0b',
  '#ec4899',
  '#06b6d4',
  '#8b5cf6',
  '#ef4444',
  '#84cc16',
  '#0ea5e9',
  '#f97316',
]

export function personColor(name: string): string {
  let hash = 0
  for (const ch of name.toLowerCase()) hash = (hash * 31 + ch.charCodeAt(0)) >>> 0
  return hues[hash % hues.length]
}

export function initials(name: string): string {
  const parts = name.trim().split(/\s+/)
  return (
    (parts[0]?.[0] ?? '') + (parts.length > 1 ? parts[parts.length - 1][0] : '')
  ).toUpperCase()
}
