import { describe, expect, it } from 'vitest'
import type { Column, EpicMap } from './api'
import { mapForSelection, progressRows } from './epicMaps'

const maps: EpicMap[] = [
  { epic: 'auth', body: '# Destination', counts: { backlog: 2, merged: 1 }, archived: 3 },
  { epic: 'billing', body: '', counts: {}, archived: 0 },
]

describe('mapForSelection', () => {
  it('returns the map of the single selected epic', () => {
    expect(mapForSelection(maps, 'epic:auth')?.epic).toBe('auth')
    expect(mapForSelection(maps, 'type:bug,epic:billing')?.epic).toBe('billing')
  })

  it('is undefined unless exactly one epic is selected', () => {
    expect(mapForSelection(maps, undefined)).toBeUndefined()
    expect(mapForSelection(maps, 'type:bug')).toBeUndefined()
    expect(mapForSelection(maps, 'epic:auth,epic:billing')).toBeUndefined()
  })

  it('is undefined for an epic without a map', () => {
    expect(mapForSelection(maps, 'epic:search')).toBeUndefined()
  })
})

describe('progressRows', () => {
  const columns = [
    { id: 'backlog', name: 'Backlog' },
    { id: 'in-review', name: 'In review' },
    { id: 'merged', name: 'Merged' },
  ] as Column[]

  it('lists every board column in board order, then Archived', () => {
    expect(progressRows(maps[0], columns)).toEqual([
      { name: 'Backlog', count: 2 },
      { name: 'In review', count: 0 },
      { name: 'Merged', count: 1 },
      { name: 'Archived', count: 3 },
    ])
  })
})
