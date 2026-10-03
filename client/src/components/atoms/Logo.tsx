import { cn } from '@/lib/utils'

export function Logo({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 32 32" aria-hidden="true" className={cn('shrink-0', className)}>
      <rect width="32" height="32" rx="8" className="fill-primary" />
      <rect x="7" y="8" width="5" height="16" rx="2.5" className="fill-primary-foreground" />
      <rect
        x="13.5"
        y="8"
        width="5"
        height="11"
        rx="2.5"
        className="fill-primary-foreground"
        opacity=".85"
      />
      <rect
        x="20"
        y="8"
        width="5"
        height="7"
        rx="2.5"
        className="fill-primary-foreground"
        opacity=".7"
      />
    </svg>
  )
}
