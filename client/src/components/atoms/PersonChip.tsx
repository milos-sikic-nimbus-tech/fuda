import { Bot } from 'lucide-react'
import { useState } from 'react'
import { Hint } from '@/components/atoms/Hint'
import { usePeople } from '@/hooks/usePeople'
import { initials, personColor } from '@/lib/colors'
import { avatarUrl } from '@/lib/people'
import { cn } from '@/lib/utils'

export function Avatar({ id, role, className }: { id: string; role?: string; className?: string }) {
  const person = usePeople().personOf(id)
  const [broken, setBroken] = useState(false)
  const src = broken ? undefined : avatarUrl(person)
  const agent = person.kind === 'agent'
  const avatar = (
    <span
      style={{ backgroundColor: personColor(person.name) }}
      className={cn(
        'inline-flex size-5 shrink-0 items-center justify-center overflow-hidden rounded-full text-[9px] font-semibold text-white',
        agent && 'rounded-md',
        className,
      )}
    >
      {agent ? (
        <Bot aria-label="Agent" className="size-[70%]" />
      ) : src ? (
        <img src={src} alt="" className="size-full object-cover" onError={() => setBroken(true)} />
      ) : (
        initials(person.name)
      )}
    </span>
  )
  if (!role) return avatar
  return <Hint label={`${role}: ${person.name}${agent ? ' (agent)' : ''}`}>{avatar}</Hint>
}

export function PersonChip({ id, onClick }: { id: string; onClick?: () => void }) {
  const person = usePeople().personOf(id)
  return (
    <button
      type="button"
      onClick={(e) => {
        e.stopPropagation()
        onClick?.()
      }}
      className="inline-flex items-center gap-1.5 rounded-full border border-border py-0.5 pr-2 pl-0.5 text-xs hover:border-primary/40"
    >
      <Avatar id={id} />
      {person.name}
      {person.kind === 'agent' && <span className="text-muted-foreground">agent</span>}
    </button>
  )
}

export function PersonName({ id }: { id: string }) {
  return usePeople().personOf(id).name
}
