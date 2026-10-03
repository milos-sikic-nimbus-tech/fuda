import { Search } from 'lucide-react'
import { useEffect, useState } from 'react'
import { Input } from '@/components/ui/input'

export function SearchInput({
  value,
  onChange,
}: {
  value?: string
  onChange: (value: string | undefined) => void
}) {
  const [draft, setDraft] = useState(value ?? '')
  const [synced, setSynced] = useState(value)
  if (value !== synced) {
    setSynced(value)
    if ((value ?? '') !== draft.trim()) setDraft(value ?? '')
  }

  useEffect(() => {
    const handle = setTimeout(() => {
      const next = draft.trim() || undefined
      if (next !== value) onChange(next)
    }, 250)
    return () => clearTimeout(handle)
  }, [draft, value, onChange])

  return (
    <div className="relative">
      <Search className="pointer-events-none absolute top-1/2 left-2.5 size-3.5 -translate-y-1/2 text-muted-foreground" />
      <Input
        value={draft}
        onChange={(e) => setDraft(e.target.value)}
        placeholder="Search id, title, people, labels"
        className="h-8 w-64 bg-card pl-8 text-sm shadow-xs"
      />
    </div>
  )
}
