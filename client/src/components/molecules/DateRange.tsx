import { format, parseISO } from 'date-fns'
import { CalendarDays, X } from 'lucide-react'
import type { DateRange as DayRange } from 'react-day-picker'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

type Range = { from?: string; to?: string }

const toDay = (value?: string) => (value ? parseISO(value) : undefined)
const toValue = (day?: Date) => (day ? format(day, 'yyyy-MM-dd') : undefined)

function describe({ from, to }: Range, emptyLabel: string): string {
  const label = (value: string) => format(parseISO(value), 'd MMM yyyy')
  if (from && to) return from === to ? label(from) : `${label(from)} – ${label(to)}`
  if (from) return `from ${label(from)}`
  if (to) return `until ${label(to)}`
  return emptyLabel
}

export function DateRange({
  from,
  to,
  onChange,
  emptyLabel = 'any time',
}: Range & { onChange: (range: Range) => void; emptyLabel?: string }) {
  const active = Boolean(from || to)
  const selected: DayRange | undefined = active ? { from: toDay(from), to: toDay(to) } : undefined

  return (
    <div className="flex items-center gap-1">
      <Popover>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            size="sm"
            className={cn(
              'h-8 flex-1 justify-start gap-2 bg-card font-normal shadow-xs',
              !active && 'text-muted-foreground',
            )}
          >
            <CalendarDays className="size-3.5" />
            {describe({ from, to }, emptyLabel)}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0" align="start">
          <Calendar
            mode="range"
            numberOfMonths={2}
            selected={selected}
            defaultMonth={selected?.from}
            onSelect={(range) => onChange({ from: toValue(range?.from), to: toValue(range?.to) })}
          />
        </PopoverContent>
      </Popover>
      {active && (
        <Button
          variant="ghost"
          size="icon"
          className="size-8"
          aria-label="Clear dates"
          onClick={() => onChange({})}
        >
          <X className="size-3.5" />
        </Button>
      )}
    </div>
  )
}
