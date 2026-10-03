import { statusHue } from '@/lib/colors'
import { cn } from '@/lib/utils'

export function StatusDot({ status, className }: { status: string; className?: string }) {
  return (
    <span
      style={{ backgroundColor: statusHue(status) }}
      className={cn('inline-block size-2 shrink-0 rounded-full', className)}
    />
  )
}
