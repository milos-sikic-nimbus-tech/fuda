import {
  Ban,
  CalendarDays,
  Check,
  ChevronLeft,
  CircleDot,
  Hand,
  Hash,
  ListFilter,
  Tag,
  TextSearch,
  User,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import { DateRange } from '@/components/molecules/DateRange'
import { ValueList } from '@/components/molecules/ValueList'
import { Button } from '@/components/ui/button'
import {
  Command,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useBoardSearch } from '@/hooks/useBoardSearch'
import type { FilterGroup } from '@/hooks/useFilterGroups'

function groupIcon(id: string) {
  if (id === 'owner' || id === 'tester') return User
  if (id === 'status') return CircleDot
  if (id === 'prefix') return Hash
  return Tag
}

function isTyping(target: EventTarget | null): boolean {
  return (
    target instanceof HTMLElement &&
    (target.isContentEditable || ['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName))
  )
}

export function FilterMenu({ groups }: { groups: FilterGroup[] }) {
  const { search, update } = useBoardSearch()
  const [open, setOpen] = useState(false)
  const [view, setView] = useState<string | null>(null)
  const current = groups.find((g) => g.id === view)

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'f' && !e.metaKey && !e.ctrlKey && !e.altKey && !isTyping(e.target)) {
        e.preventDefault()
        setOpen(true)
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [])

  const flags = [
    {
      id: 'blocked',
      label: 'Blocked only',
      icon: Ban,
      on: !!search.blocked,
      toggle: () => update({ blocked: search.blocked ? undefined : true }),
    },
    {
      id: 'available',
      label: 'Available to pick up',
      icon: Hand,
      on: !!search.available,
      toggle: () => update({ available: search.available ? undefined : true }),
    },
    {
      id: 'text',
      label: 'Search also in task text',
      icon: TextSearch,
      on: !!search.text,
      toggle: () => update({ text: search.text ? undefined : true }),
    },
  ]

  return (
    <Popover
      open={open}
      onOpenChange={(next) => {
        setOpen(next)
        if (!next) setView(null)
      }}
    >
      <PopoverTrigger asChild>
        <Button variant="outline" size="sm" className="h-8 gap-1.5 bg-card shadow-xs">
          <ListFilter className="size-3.5" />
          Filter
          <kbd className="ml-1 hidden rounded border bg-muted px-1 font-mono text-[10px] text-muted-foreground sm:inline">
            F
          </kbd>
        </Button>
      </PopoverTrigger>
      <PopoverContent align="start" className="w-[min(18rem,calc(100vw-1.5rem))] p-0">
        {view === 'added' ? (
          <div className="p-3">
            <BackButton onClick={() => setView(null)} label="Added" />
            <DateRange
              from={search.from}
              to={search.to}
              onChange={(r) => update({ from: r.from, to: r.to })}
            />
          </div>
        ) : current ? (
          <Command key={current.id}>
            <div className="border-b p-1">
              <BackButton onClick={() => setView(null)} label={current.title} />
            </div>
            <ValueList group={current} autoFocus />
          </Command>
        ) : (
          <Command key="groups">
            <CommandInput autoFocus placeholder="Filter by…" />
            <CommandList className="max-h-96">
              <CommandGroup>
                {groups.map((g) => {
                  const Icon = groupIcon(g.id)
                  return (
                    <CommandItem
                      key={g.id}
                      value={g.title}
                      onSelect={() => setView(g.id)}
                      className="gap-2.5"
                    >
                      <Icon className="size-3.5 text-muted-foreground" />
                      <span className="flex-1">{g.title}</span>
                      {g.selected.length > 0 && (
                        <span className="text-xs text-primary tabular-nums">
                          {g.selected.length}
                        </span>
                      )}
                    </CommandItem>
                  )
                })}
                <CommandItem
                  value="Added date"
                  onSelect={() => setView('added')}
                  className="gap-2.5"
                >
                  <CalendarDays className="size-3.5 text-muted-foreground" />
                  <span className="flex-1">Added</span>
                  {(search.from || search.to) && (
                    <span className="size-1.5 rounded-full bg-primary" />
                  )}
                </CommandItem>
              </CommandGroup>
              <CommandSeparator />
              <CommandGroup>
                {flags.map(({ id, label, icon: Icon, on, toggle }) => (
                  <CommandItem key={id} value={label} onSelect={toggle} className="gap-2.5">
                    <Icon className="size-3.5 text-muted-foreground" />
                    <span className="flex-1">{label}</span>
                    {on && <Check className="size-3.5 text-primary" />}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        )}
      </PopoverContent>
    </Popover>
  )
}

function BackButton({ label, onClick }: { label: string; onClick: () => void }) {
  return (
    <button
      type="button"
      onClick={onClick}
      className="mb-1 flex w-full items-center gap-1 rounded-md px-1.5 py-1 text-xs font-medium text-muted-foreground hover:bg-accent hover:text-foreground"
    >
      <ChevronLeft className="size-3.5" />
      {label}
    </button>
  )
}
