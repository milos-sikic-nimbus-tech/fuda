import { useCallback, useMemo } from 'react'
import type { Person } from '@/lib/api'
import { personOf } from '@/lib/people'
import { useBoard } from '@/lib/queries'

export function usePeople(): { people: Person[]; personOf: (id: string) => Person } {
  const { data } = useBoard()
  const people = useMemo(() => data?.facets.people ?? [], [data])
  return { people, personOf: useCallback((id: string) => personOf(people, id), [people]) }
}
