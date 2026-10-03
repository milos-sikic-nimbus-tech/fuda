import { Input } from '@/components/ui/input'

export function DateRange({
  from,
  to,
  onChange,
}: {
  from?: string
  to?: string
  onChange: (range: { from?: string; to?: string }) => void
}) {
  return (
    <div className="flex items-center gap-1 text-xs text-muted-foreground">
      <span>Added</span>
      <Input
        type="date"
        value={from ?? ''}
        onChange={(e) => onChange({ from: e.target.value || undefined, to })}
        className="h-8 w-34 bg-card text-xs"
      />
      <span>–</span>
      <Input
        type="date"
        value={to ?? ''}
        onChange={(e) => onChange({ from, to: e.target.value || undefined })}
        className="h-8 w-34 bg-card text-xs"
      />
    </div>
  )
}
