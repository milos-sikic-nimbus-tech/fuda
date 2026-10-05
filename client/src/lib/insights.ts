import {
  addDays,
  differenceInCalendarDays,
  format,
  isValid,
  parseISO,
  startOfDay,
  startOfISOWeek,
  subDays,
  subWeeks,
} from 'date-fns'
import { z } from 'zod'
import type { Card } from './api'
import { labelText } from '@/lib/labels'

export const rangeKeys = ['7d', '4w', '12w', '26w', '52w'] as const
export type RangeKey = (typeof rangeKeys)[number]

const isoDate = z
  .string()
  .regex(/^\d{4}-\d{2}-\d{2}$/)
  .refine((value) => isValid(parseISO(value)))
  .optional()
  .catch(undefined)

export const insightsSearchSchema = z.object({
  range: z.enum(rangeKeys).optional().catch(undefined),
  from: isoDate,
  to: isoDate,
})

export type InsightsSearch = z.infer<typeof insightsSearchSchema>

export type Period = { start: Date; end: Date; label: string }

const dayLabel = (day: Date) => format(day, 'd MMM yyyy')

export function periodOf({ range = '12w', from, to }: InsightsSearch, today = new Date()): Period {
  const end = startOfDay(today)
  if (from) {
    const start = parseISO(from)
    const last = to ? parseISO(to) : end
    return {
      start,
      end: last,
      label: from === to ? dayLabel(start) : `${dayLabel(start)} – ${dayLabel(last)}`,
    }
  }
  if (range === '7d') return { start: subDays(end, 6), end, label: 'last 7 days' }
  const weeks = Number.parseInt(range)
  return { start: subWeeks(startOfISOWeek(end), weeks - 1), end, label: `last ${weeks} weeks` }
}

export function within(period: Period, date: string): boolean {
  if (!date) return false
  const day = parseISO(date)
  return day >= period.start && day <= period.end
}

export type BucketUnit = 'day' | 'week'
export type Bucket = { start: string; label: string; created: number; claimed: number }

const maxDailyBuckets = 14

export function activity(cards: Card[], period: Period): { unit: BucketUnit; buckets: Bucket[] } {
  const days = differenceInCalendarDays(period.end, period.start) + 1
  const unit: BucketUnit = days <= maxDailyBuckets ? 'day' : 'week'
  const step = unit === 'day' ? 1 : 7
  const first = unit === 'day' ? period.start : startOfISOWeek(period.start)
  const count = Math.max(0, Math.floor(differenceInCalendarDays(period.end, first) / step) + 1)
  const buckets: Bucket[] = Array.from({ length: count }, (_, i) => {
    const start = addDays(first, i * step)
    return {
      start: format(start, 'yyyy-MM-dd'),
      label: format(start, 'd MMM'),
      created: 0,
      claimed: 0,
    }
  })
  const place = (date: string, key: 'created' | 'claimed') => {
    if (!within(period, date)) return
    buckets[Math.floor(differenceInCalendarDays(parseISO(date), first) / step)][key]++
  }
  for (const card of cards) {
    place(card.added, 'created')
    place(card.claimed, 'claimed')
  }
  return { unit, buckets }
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
