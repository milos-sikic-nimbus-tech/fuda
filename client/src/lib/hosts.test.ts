import { describe, expect, it } from 'vitest'
import { loginRedirect, onlyBoard, type BoardListData } from './hosts'

const github = { host: 'github', loggedIn: true, login: '/auth/github/login' } as const
const azure = { host: 'azure', loggedIn: false, login: '/auth/azure/login' } as const
const board = {
  host: 'github',
  repo: 'o/fuda-x',
  path: '/github/o/fuda-x',
  title: 'fuda-x',
} as const

describe('loginRedirect', () => {
  it('sends a person with one logged-out host straight to its login', () => {
    const listing: BoardListData = { boards: [], hosts: [{ ...github, loggedIn: false }] }
    expect(loginRedirect(listing)).toBe('/auth/github/login')
  })
  it('lets a person with two hosts choose', () => {
    const listing: BoardListData = {
      boards: [],
      hosts: [{ ...github, loggedIn: false }, azure],
    }
    expect(loginRedirect(listing)).toBeUndefined()
  })
  it('stays put when a host is logged in', () => {
    expect(loginRedirect({ boards: [], hosts: [github] })).toBeUndefined()
  })
})

describe('onlyBoard', () => {
  it('opens a single Board when every host is logged in', () => {
    expect(onlyBoard({ boards: [board], hosts: [github] })).toEqual(board)
  })
  it('keeps the list when another host still needs a login', () => {
    expect(onlyBoard({ boards: [board], hosts: [github, azure] })).toBeUndefined()
  })
})
