import type { Column, EpicMap } from './api'
import { listOf } from './filters'

export function mapForSelection(
  maps: EpicMap[],
  labelParam: string | undefined,
): EpicMap | undefined {
  const epics = listOf(labelParam)
    .filter((l) => l.startsWith('epic:'))
    .map((l) => l.slice('epic:'.length))
  if (epics.length !== 1) return undefined
  return maps.find((m) => m.epic === epics[0])
}

export function progressRows(map: EpicMap, columns: Column[]): { name: string; count: number }[] {
  return [
    ...columns.map((c) => ({ name: c.name, count: map.counts[c.id] ?? 0 })),
    { name: 'Archived', count: map.archived },
  ]
}
