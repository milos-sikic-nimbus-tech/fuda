import { describe, expect, it } from 'vitest'
import { boardPrefix } from './boardPath'

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
