const knownStatuses = [
  'backlog',
  'in progress',
  'in review',
  'merged',
  'testing',
  'validated',
  'done',
]
const knownTypes = ['bug', 'feat', 'impr', 'refactor', 'test', 'chore', 'docs', 'question']

function token(prefix: string, known: string[], value: string): string {
  const key = value.trim().toLowerCase()
  return `var(--color-${prefix}-${known.includes(key) ? key.replaceAll(' ', '-') : 'other'})`
}

export const statusHue = (status: string) => token('status', knownStatuses, status)

export const typeHue = (type: string) => token('type', knownTypes, type)

const personHues = [
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
  let h = 0
  for (const ch of name.toLowerCase()) h = (h * 31 + ch.charCodeAt(0)) >>> 0
  return personHues[h % personHues.length]
}

export function initials(name: string): string {
  const parts = name.trim().split(/\s+/)
  return (
    (parts[0]?.[0] ?? '') + (parts.length > 1 ? parts[parts.length - 1][0] : '')
  ).toUpperCase()
}
