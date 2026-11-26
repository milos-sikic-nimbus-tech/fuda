import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import type { Column, EpicMap } from '@/lib/api'
import { progressRows } from '@/lib/epicMaps'

export function EpicMapDialog({
  map,
  columns,
  open,
  onOpenChange,
}: {
  map: EpicMap
  columns: Column[]
  open: boolean
  onOpenChange: (open: boolean) => void
}) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[85vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>Map: {map.epic}</DialogTitle>
          <DialogDescription>Tasks of this epic per column.</DialogDescription>
        </DialogHeader>
        <dl className="flex flex-wrap gap-2">
          {progressRows(map, columns).map((row) => (
            <div key={row.name} className="rounded-lg border border-border px-3 py-1.5">
              <dt className="text-xs text-muted-foreground">{row.name}</dt>
              <dd className="text-lg font-semibold tabular-nums">{row.count}</dd>
            </div>
          ))}
        </dl>
        <pre className="font-mono text-xs break-words whitespace-pre-wrap">{map.body}</pre>
      </DialogContent>
    </Dialog>
  )
}
