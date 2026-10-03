import { describe, expect, it } from 'vitest'
import type { Card } from './api'
import { countBy, isOpen, labelsOf, topSlices, weeklyActivity } from './insights'

function card(id: string, overrides: Partial<Card> = {}): Card {
  return {
    id,
    title: id,
    status: 'backlog',
    column: 'backlog',
    owners: [],
    testers: [],
    labels: [],
    prefix: 'X',
    added: '',
    claimed: '',
    prs: [],
    openPrs: [],
    blockedByIds: [],
    blocks: [],
    references: [],
    referencedBy: [],
    inProd: false,
    newInPr: false,
    archiveCandidate: false,
    path: '',
    ...overrides,
  }
}

describe('weeklyActivity', () => {
  const today = new Date('2026-10-03T12:00:00')
  it('buckets created and claimed dates into ISO weeks ending this week', () => {
    const weeks = weeklyActivity(
      [
        card('A', { added: '2026-09-29', claimed: '2026-10-02' }),
        card('B', { added: '2026-09-22' }),
        card('C', { added: '2026-08-01' }),
      ],
      2,
      today,
    )
    expect(weeks.map((w) => w.start)).toEqual(['2026-09-21', '2026-09-28'])
    expect(weeks.map((w) => [w.created, w.claimed])).toEqual([
      [1, 0],
      [1, 1],
    ])
  })
})

describe('breakdowns', () => {
  const cards = [
    card('A', { labels: ['epic:billing', 'type:bug'] }),
    card('B', { labels: ['epic:billing'] }),
    card('C', { labels: ['type:feat'], status: 'merged' }),
  ]
  it('counts by label group with None for unlabeled', () => {
    expect(countBy(cards, labelsOf('epic'))).toEqual(
      new Map([
        ['billing', 2],
        ['None', 1],
      ]),
    )
  })
  it('keeps the top slices and folds the rest into Other', () => {
    const counts = new Map([
      ['a', 5],
      ['b', 4],
      ['c', 3],
      ['d', 1],
    ])
    expect(topSlices(counts, 3)).toEqual([
      { key: 'a', count: 5 },
      { key: 'b', count: 4 },
      { key: 'Other', count: 4 },
    ])
  })
  it('treats merged and later as not open', () => {
    expect(cards.filter(isOpen).map((c) => c.id)).toEqual(['A', 'B'])
  })
})
