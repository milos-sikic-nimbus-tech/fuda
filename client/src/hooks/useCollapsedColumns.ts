import { useCallback, useState } from 'react'

const storageKey = 'fuda-collapsed-columns'

function read(): string[] {
  try {
    const value = JSON.parse(localStorage.getItem(storageKey) ?? '[]')
    return Array.isArray(value) ? value.filter((v) => typeof v === 'string') : []
  } catch {
    return []
  }
}

export function useCollapsedColumns() {
  const [collapsed, setCollapsed] = useState<string[]>(read)

  const toggle = useCallback((id: string) => {
    setCollapsed((current) => {
      const next = current.includes(id) ? current.filter((c) => c !== id) : [...current, id]
      try {
        localStorage.setItem(storageKey, JSON.stringify(next))
      } catch {
        return next
      }
      return next
    })
  }, [])

  return { isCollapsed: (id: string) => collapsed.includes(id), toggle }
}
