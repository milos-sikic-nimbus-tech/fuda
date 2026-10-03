import { Checkbox } from '@/components/ui/checkbox'
import {
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import type { FilterGroup } from '@/hooks/useFilterGroups'

export function ValueList({ group, autoFocus }: { group: FilterGroup; autoFocus?: boolean }) {
  return (
    <>
      <CommandInput autoFocus={autoFocus} placeholder={`Filter ${group.plural}…`} />
      <CommandList className="max-h-80">
        <CommandEmpty>No {group.plural} match.</CommandEmpty>
        <CommandGroup>
          {group.values.map((v) => (
            <CommandItem
              key={v.value}
              value={v.text}
              onSelect={() => group.toggle(v.value)}
              className="gap-2.5"
            >
              <Checkbox
                checked={group.selected.includes(v.value)}
                className="pointer-events-none"
                tabIndex={-1}
              />
              <span className="min-w-0 flex-1 truncate">{v.label}</span>
              <span className="text-xs text-muted-foreground tabular-nums">{v.count}</span>
            </CommandItem>
          ))}
        </CommandGroup>
      </CommandList>
    </>
  )
}
