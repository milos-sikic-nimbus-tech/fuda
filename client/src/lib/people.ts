import type { Person } from './api'

export const kindValue = (kind: Person['kind']) => `kind:${kind}`

export function personOf(people: Person[], id: string): Person {
  return people.find((p) => p.id === id) ?? { id, name: id, aliases: [], kind: 'human' }
}

export function avatarUrl(person: Person): string | undefined {
  return person.github
    ? `https://github.com/${encodeURIComponent(person.github)}.png?size=64`
    : undefined
}

export function ownerMatches(people: Person[], selected: string, owners: string[]): boolean {
  if (selected.startsWith('kind:')) {
    return owners.some((o) => kindValue(personOf(people, o).kind) === selected)
  }
  return owners.includes(selected)
}
