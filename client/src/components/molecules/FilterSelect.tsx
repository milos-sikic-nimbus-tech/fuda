import { ChevronDown } from 'lucide-react'
import type { ReactNode } from 'react'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

export type Option = { value: string; label: ReactNode; search?: string; count?: number }

const searchableFrom = 8

function summary(title: string, options: Option[], selected: string[]): ReactNode {
  if (selected.length === 0) return title
  if (selected.length > 1) return `${title} · ${selected.length}`
  const only = options.find((o) => o.value === selected[0])
  return (
    <>
      <span className="text-muted-foreground">{title}:</span>
      <span className="max-w-32 truncate">{only?.search ?? selected[0]}</span>
    </>
  )
}

export function FilterSelect({
  title,
  options,
  selected,
  onToggle,
  onClear,
}: {
  title: string
  options: Option[]
  selected: string[]
  onToggle: (value: string) => void
  onClear?: () => void
}) {
  if (options.length === 0) return null
  const active = selected.length > 0

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className={cn(
            'h-8 gap-1.5 bg-card shadow-xs',
            active && 'border-primary/50 bg-primary/5 text-primary',
          )}
        >
          {summary(title, options, selected)}
          <ChevronDown className="size-3.5 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-64 p-0" align="start">
        <Command>
          {options.length >= searchableFrom && (
            <CommandInput placeholder={`Find ${title.toLowerCase()}…`} />
          )}
          <CommandList className="max-h-72">
            <CommandEmpty>Nothing matches.</CommandEmpty>
            <CommandGroup>
              {options.map((o) => (
                <CommandItem
                  key={o.value}
                  value={o.search ?? o.value}
                  onSelect={() => onToggle(o.value)}
                  className="gap-2.5"
                >
                  <Checkbox
                    checked={selected.includes(o.value)}
                    className="pointer-events-none"
                    tabIndex={-1}
                  />
                  <span className="min-w-0 flex-1 truncate">{o.label}</span>
                  {o.count !== undefined && (
                    <span className="text-xs text-muted-foreground tabular-nums">{o.count}</span>
                  )}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
          {active && onClear && (
            <div className="border-t p-1">
              <Button variant="ghost" size="sm" className="h-7 w-full text-xs" onClick={onClear}>
                Clear {title.toLowerCase()}
              </Button>
            </div>
          )}
        </Command>
      </PopoverContent>
    </Popover>
  )
}
