import { Check, UserPlus } from 'lucide-react'
import { useState } from 'react'
import { Avatar } from '@/components/atoms/PersonChip'
import { Button } from '@/components/ui/button'
import {
  Command,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import type { Person } from '@/lib/api'
import { ownerChoices, toggleOwner } from '@/lib/assigns'
import { personOf } from '@/lib/people'

export function OwnerPicker({
  people,
  owners,
  me,
  saving,
  readOnly,
  onChange,
}: {
  people: Person[]
  owners: string[]
  me?: string
  saving: boolean
  readOnly: boolean
  onChange: (owners: string[]) => void
}) {
  const [open, setOpen] = useState(false)

  return (
    <>
      {me && !owners.includes(me) && (
        <Button
          variant="ghost"
          size="xs"
          disabled={saving || readOnly}
          className="text-muted-foreground"
          onClick={() => onChange([...owners, me])}
        >
          Assign to me
        </Button>
      )}
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="ghost"
            size="xs"
            disabled={saving || readOnly}
            aria-label="Pick owners"
            className="text-muted-foreground"
          >
            <UserPlus className="size-3.5" />
            {saving && 'Saving…'}
          </Button>
        </PopoverTrigger>
        <PopoverContent align="start" className="w-56 p-0">
          <Command>
            <CommandInput placeholder="Find a person…" />
            <CommandList>
              <CommandEmpty>No one found.</CommandEmpty>
              {ownerChoices(
                people.map((p) => p.id),
                owners,
              ).map((id) => {
                const person = personOf(people, id)
                return (
                  <CommandItem
                    key={id}
                    value={id}
                    keywords={[
                      person.name,
                      ...person.aliases,
                      ...(person.github ? [person.github] : []),
                    ]}
                    onSelect={() => {
                      setOpen(false)
                      onChange(toggleOwner(owners, id))
                    }}
                  >
                    <Avatar id={id} className="size-4.5 text-[8px]" />
                    <span className="min-w-0 flex-1 truncate">{person.name}</span>
                    {owners.includes(id) && <Check className="size-3.5" />}
                  </CommandItem>
                )
              })}
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </>
  )
}
