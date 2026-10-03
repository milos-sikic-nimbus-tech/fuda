import { SearchX } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useBoardSearch } from '@/hooks/useBoardSearch'

export function NoMatches() {
  const { clear } = useBoardSearch()
  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-3 p-8 text-center">
      <div className="flex size-10 items-center justify-center rounded-full bg-muted text-muted-foreground">
        <SearchX className="size-5" />
      </div>
      <div>
        <p className="text-sm font-medium">No tasks match these filters</p>
        <p className="mt-1 text-xs text-muted-foreground">Remove a filter above, or start over.</p>
      </div>
      <Button variant="outline" size="sm" onClick={clear}>
        Clear filters
      </Button>
    </div>
  )
}
