import { describe, expect, it } from 'vitest'
import { buildDocTree, foldersOf } from './docTree'

describe('buildDocTree', () => {
  it('puts folders first, then files, each A to Z', () => {
    const tree = buildDocTree([
      'zeta.md',
      'board/stages.md',
      'alpha.md',
      'board/README.md',
      'adr/0001.md',
    ])
    expect(tree.map((n) => n.name)).toEqual(['adr', 'board', 'alpha.md', 'zeta.md'])
    const board = tree[1]
    expect(board.kind === 'folder' && board.children.map((n) => n.path)).toEqual([
      'board/README.md',
      'board/stages.md',
    ])
  })

  it('nests folders', () => {
    const [guides] = buildDocTree(['guides/deep/one.md'])
    expect(guides.kind === 'folder' && guides.children[0]).toMatchObject({
      kind: 'folder',
      path: 'guides/deep',
    })
  })

  it('is empty without docs', () => {
    expect(buildDocTree([])).toEqual([])
  })
})

describe('foldersOf', () => {
  it('lists every folder above a doc, outermost first', () => {
    expect(foldersOf('a/b/c.md')).toEqual(['a', 'a/b'])
    expect(foldersOf('c.md')).toEqual([])
  })
})
