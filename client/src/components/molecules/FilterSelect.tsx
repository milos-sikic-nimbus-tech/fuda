import { Check, ChevronDown } from 'lucide-react'
import type { ReactNode } from 'react'
import { Button } from '@/components/ui/button'
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

export type Option = { value: string; label: ReactNode; search?: string }

export function FilterSelect({
  title,
  options,
  selected,
  onToggle,
}: {
  title: string
  options: Option[]
  selected: string[]
  onToggle: (value: string) => void
}) {
  if (options.length === 0) return null
  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className={cn('h-8 gap-1.5', selected.length > 0 && 'border-primary text-primary')}
        >
          {title}
          {selected.length > 0 && (
            <span className="rounded bg-primary px-1 text-[10px] font-semibold text-primary-foreground">
              {selected.length}
            </span>
          )}
          <ChevronDown className="size-3.5 opacity-60" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-64 p-0" align="start">
        <Command>
          {options.length > 7 && <CommandInput placeholder={`Filter ${title.toLowerCase()}…`} />}
          <CommandList>
            <CommandEmpty>Nothing matches.</CommandEmpty>
            <CommandGroup>
              {options.map((o) => (
                <CommandItem
                  key={o.value}
                  value={o.search ?? o.value}
                  onSelect={() => onToggle(o.value)}
                >
                  <span
                    className={cn(
                      'flex size-4 items-center justify-center rounded border',
                      selected.includes(o.value)
                        ? 'border-primary bg-primary text-primary-foreground'
                        : 'border-input',
                    )}
                  >
                    {selected.includes(o.value) && <Check className="size-3" />}
                  </span>
                  <span className="truncate">{o.label}</span>
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
