import {
  addDays,
  differenceInCalendarDays,
  format,
  parseISO,
  startOfISOWeek,
  subWeeks,
} from 'date-fns'
import { z } from 'zod'
import type { Card } from './api'
import { labelText } from '@/lib/labels'

export const ranges = { '4w': 4, '12w': 12, '26w': 26, '52w': 52 } as const
export type RangeKey = keyof typeof ranges

export const insightsSearchSchema = z.object({
  range: z.enum(['4w', '12w', '26w', '52w']).optional().catch(undefined),
})

export type WeekBucket = { start: string; label: string; created: number; claimed: number }

export function weeklyActivity(cards: Card[], weeks: number, today = new Date()): WeekBucket[] {
  const first = subWeeks(startOfISOWeek(today), weeks - 1)
  const buckets: WeekBucket[] = Array.from({ length: weeks }, (_, i) => {
    const start = addDays(first, i * 7)
    return {
      start: format(start, 'yyyy-MM-dd'),
      label: format(start, 'd MMM'),
      created: 0,
      claimed: 0,
    }
  })
  const place = (date: string, key: 'created' | 'claimed') => {
    if (!date) return
    const offset = differenceInCalendarDays(parseISO(date), first)
    if (offset < 0) return
    const index = Math.floor(offset / 7)
    if (index < buckets.length) buckets[index][key]++
  }
  for (const card of cards) {
    place(card.added, 'created')
    place(card.claimed, 'claimed')
  }
  return buckets
}

export type Slice = { key: string; count: number }

export function topSlices(counts: Map<string, number>, limit = 8): Slice[] {
  const sorted = [...counts.entries()]
    .map(([key, count]) => ({ key, count }))
    .sort((a, b) => b.count - a.count || a.key.localeCompare(b.key))
  if (sorted.length <= limit) return sorted
  const rest = sorted.slice(limit - 1).reduce((n, s) => n + s.count, 0)
  return [...sorted.slice(0, limit - 1), { key: 'Other', count: rest }]
}

export function countBy(cards: Card[], keysOf: (card: Card) => string[]): Map<string, number> {
  const counts = new Map<string, number>()
  for (const card of cards) {
    const keys = keysOf(card)
    for (const key of keys.length ? keys : ['None']) counts.set(key, (counts.get(key) ?? 0) + 1)
  }
  return counts
}

export const labelsOf = (group: string) => (card: Card) =>
  card.labels.filter((l) => labelText(l).group === group).map((l) => labelText(l).value)

export function isOpen(card: Card): boolean {
  return !['merged', 'testing', 'validated', 'done'].includes(card.status.trim().toLowerCase())
}
