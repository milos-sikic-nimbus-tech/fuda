import { Link } from '@tanstack/react-router'
import { ChevronRight, FileText } from 'lucide-react'
import { useState } from 'react'
import { buildDocTree, type DocNode } from '@/lib/docTree'
import { cn } from '@/lib/utils'

export function DocTree({ paths, openPath }: { paths: string[]; openPath: string }) {
  const [collapsed, setCollapsed] = useState<string[]>([])

  const toggle = (path: string) =>
    setCollapsed((current) =>
      current.includes(path) ? current.filter((p) => p !== path) : [...current, path],
    )
  const isExpanded = (path: string) => !collapsed.includes(path)

  const renderNode = (node: DocNode, depth: number) => {
    const indent = { paddingLeft: `${depth * 12 + 8}px` }
    if (node.kind === 'doc') {
      return (
        <Link
          key={node.path}
          to="/docs/$"
          params={{ _splat: node.path }}
          style={indent}
          aria-current={node.path === openPath ? 'page' : undefined}
          className={cn(
            'flex items-center gap-1.5 rounded-md py-1 pr-2 text-sm text-muted-foreground hover:bg-accent hover:text-foreground',
            node.path === openPath && 'bg-accent font-medium text-foreground',
          )}
        >
          <FileText className="size-3.5 shrink-0" />
          <span className="truncate">{node.name}</span>
        </Link>
      )
    }
    const expanded = isExpanded(node.path)
    return (
      <div key={node.path}>
        <button
          type="button"
          style={indent}
          aria-expanded={expanded}
          onClick={() => toggle(node.path)}
          className="flex w-full items-center gap-1.5 rounded-md py-1 pr-2 text-sm text-muted-foreground hover:bg-accent hover:text-foreground"
        >
          <ChevronRight
            className={cn('size-3.5 shrink-0 transition-transform', expanded && 'rotate-90')}
          />
          <span className="truncate">{node.name}</span>
        </button>
        {expanded && node.children.map((child) => renderNode(child, depth + 1))}
      </div>
    )
  }

  return (
    <nav aria-label="Docs" className="w-60 shrink-0 overflow-y-auto border-r border-border/70 p-2">
      {buildDocTree(paths).map((node) => renderNode(node, 0))}
    </nav>
  )
}
