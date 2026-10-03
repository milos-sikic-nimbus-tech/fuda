import { Table2 } from 'lucide-react'
import { type ReactNode, useState } from 'react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

export function ChartCard({
  title,
  description,
  legend,
  table,
  children,
  className,
}: {
  title: string
  description?: string
  legend?: ReactNode
  table?: ReactNode
  children: ReactNode
  className?: string
}) {
  const [asTable, setAsTable] = useState(false)
  return (
    <section className={cn('rounded-xl border border-border bg-card p-4 shadow-xs', className)}>
      <header className="mb-3 flex flex-wrap items-start gap-x-4 gap-y-2">
        <div className="min-w-0 flex-1">
          <h2 className="text-sm font-medium">{title}</h2>
          {description && <p className="mt-0.5 text-xs text-muted-foreground">{description}</p>}
        </div>
        {legend}
        {table && (
          <Button
            variant="ghost"
            size="sm"
            className="h-7 gap-1.5 px-2 text-xs text-muted-foreground"
            aria-pressed={asTable}
            onClick={() => setAsTable((v) => !v)}
          >
            <Table2 className="size-3.5" />
            {asTable ? 'Chart' : 'Table'}
          </Button>
        )}
      </header>
      {asTable && table ? table : children}
    </section>
  )
}

export function LegendItem({ color, label }: { color: string; label: string }) {
  return (
    <span className="inline-flex items-center gap-1.5 text-xs text-muted-foreground">
      <span className="size-2.5 rounded-sm" style={{ backgroundColor: color }} />
      {label}
    </span>
  )
}
