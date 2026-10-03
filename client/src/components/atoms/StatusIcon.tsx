import { statusHue } from '@/lib/colors'
import { cn } from '@/lib/utils'

type Shape = 'dashed' | 'empty' | 'quarter' | 'half' | 'three-quarter' | 'check'

const shapes: Record<string, Shape> = {
  backlog: 'dashed',
  'in progress': 'half',
  'in review': 'three-quarter',
  merged: 'check',
  testing: 'quarter',
  validated: 'check',
  done: 'check',
}

function arc(fraction: number): string {
  const angle = fraction * 2 * Math.PI
  const x = 7 + 4 * Math.sin(angle)
  const y = 7 - 4 * Math.cos(angle)
  return `M7 7 L7 3 A4 4 0 ${fraction > 0.5 ? 1 : 0} 1 ${x.toFixed(2)} ${y.toFixed(2)} Z`
}

export function StatusIcon({ status, className }: { status: string; className?: string }) {
  const shape = shapes[status.trim().toLowerCase()] ?? 'empty'
  const color = statusHue(status)
  return (
    <svg
      viewBox="0 0 14 14"
      aria-hidden="true"
      className={cn('size-3.5 shrink-0', className)}
      style={{ color }}
    >
      {shape === 'check' ? (
        <>
          <circle cx="7" cy="7" r="6" fill="currentColor" />
          <path
            d="M4.4 7.2 6.2 9l3.4-3.6"
            fill="none"
            stroke="var(--color-card)"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </>
      ) : (
        <>
          <circle
            cx="7"
            cy="7"
            r="6"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeDasharray={shape === 'dashed' ? '2 2.2' : undefined}
          />
          {shape === 'quarter' && <path d={arc(0.25)} fill="currentColor" />}
          {shape === 'half' && <path d={arc(0.5)} fill="currentColor" />}
          {shape === 'three-quarter' && <path d={arc(0.75)} fill="currentColor" />}
        </>
      )}
    </svg>
  )
}
