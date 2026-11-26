import { describe, expect, it } from 'vitest'
import type { Person } from './api'
import { avatarUrl, ownerMatches, personOf } from './people'

const people: Person[] = [
  { id: 'ann', name: 'Ann Lee', aliases: ['Annie'], github: 'ann-gh', kind: 'human' },
  { id: 'research-agent', name: 'Research agent', aliases: [], kind: 'agent' },
]

describe('personOf', () => {
  it('finds a listed person by id', () => {
    expect(personOf(people, 'ann').name).toBe('Ann Lee')
  })

  it('shows an unknown owner as itself', () => {
    expect(personOf(people, 'Guest')).toEqual({
      id: 'Guest',
      name: 'Guest',
      aliases: [],
      kind: 'human',
    })
  })
})

describe('avatarUrl', () => {
  it('uses the GitHub avatar when a login is set', () => {
    expect(avatarUrl(people[0])).toBe('https://github.com/ann-gh.png?size=64')
    expect(avatarUrl(people[1])).toBeUndefined()
  })
})

describe('ownerMatches', () => {
  it('matches ids and the human or agent split', () => {
    expect(ownerMatches(people, 'ann', ['ann'])).toBe(true)
    expect(ownerMatches(people, 'kind:agent', ['research-agent'])).toBe(true)
    expect(ownerMatches(people, 'kind:agent', ['ann'])).toBe(false)
    expect(ownerMatches(people, 'kind:human', ['ann', 'research-agent'])).toBe(true)
    expect(ownerMatches(people, 'kind:human', ['Guest'])).toBe(true)
  })
})
