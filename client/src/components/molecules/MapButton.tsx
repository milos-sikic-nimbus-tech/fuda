import { Map as MapIcon } from 'lucide-react'
import { useState } from 'react'
import { Hint } from '@/components/atoms/Hint'
import { EpicMapDialog } from '@/components/organisms/EpicMapDialog'
import { Button } from '@/components/ui/button'
import { useBoardSearch } from '@/hooks/useBoardSearch'
import type { Column, EpicMap } from '@/lib/api'
import { mapForSelection } from '@/lib/epicMaps'

export function MapButton({ maps, columns }: { maps: EpicMap[]; columns: Column[] }) {
  const { search } = useBoardSearch()
  const [open, setOpen] = useState(false)
  if (maps.length === 0) return null
  const selected = mapForSelection(maps, search.label)

  return (
    <>
      <Hint
        label={
          selected ? `Show the map of ${selected.epic}` : 'Select one epic with a map in the filter'
        }
      >
        <span>
          <Button
            variant="outline"
            size="sm"
            className="h-8 gap-1.5 bg-card shadow-xs"
            disabled={!selected}
            onClick={() => setOpen(true)}
            aria-label="Epic map"
          >
            <MapIcon className="size-3.5" />
            Map
          </Button>
        </span>
      </Hint>
      {selected && (
        <EpicMapDialog map={selected} columns={columns} open={open} onOpenChange={setOpen} />
      )}
    </>
  )
}
