import { describe, expect, it } from 'vitest'
import { boardPrefix, loginPath } from './boardPath'

describe('boardPrefix', () => {
  it('takes the host and its repository segments', () => {
    expect(boardPrefix('/github/o/r/board')).toBe('/github/o/r')
    expect(boardPrefix('/azure/org/proj/repo/archive')).toBe('/azure/org/proj/repo')
    expect(boardPrefix('/local/folder/docs/a.md')).toBe('/local/folder')
  })
  it('is empty outside a Board', () => {
    expect(boardPrefix('/')).toBe('')
    expect(boardPrefix('/guide/setup')).toBe('')
    expect(boardPrefix('/github/o')).toBe('')
  })
})

describe('loginPath', () => {
  it('comes back to the page the person was on', () => {
    expect(loginPath('/github/o/fuda-x/archive?q=a b')).toBe(
      '/auth/github/login?return=%2Fgithub%2Fo%2Ffuda-x%2Farchive%3Fq%3Da%20b',
    )
  })
})
