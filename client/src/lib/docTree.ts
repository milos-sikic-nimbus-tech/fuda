export type DocNode =
  | { kind: 'folder'; name: string; path: string; children: DocNode[] }
  | { kind: 'doc'; name: string; path: string }

export function buildDocTree(paths: string[]): DocNode[] {
  const root: DocNode[] = []
  for (const path of paths) {
    const parts = path.split('/')
    let level = root
    parts.forEach((name, i) => {
      const nodePath = parts.slice(0, i + 1).join('/')
      if (i === parts.length - 1) {
        level.push({ kind: 'doc', name, path: nodePath })
        return
      }
      let folder = level.find((n) => n.kind === 'folder' && n.name === name)
      if (!folder) {
        folder = { kind: 'folder', name, path: nodePath, children: [] }
        level.push(folder)
      }
      level = folder.kind === 'folder' ? folder.children : level
    })
  }
  return sortNodes(root)
}

function sortNodes(nodes: DocNode[]): DocNode[] {
  return nodes
    .map((n) => (n.kind === 'folder' ? { ...n, children: sortNodes(n.children) } : n))
    .sort((a, b) =>
      a.kind === b.kind ? a.name.localeCompare(b.name) : a.kind === 'folder' ? -1 : 1,
    )
}
