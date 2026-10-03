import { FileWarning } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'
import type { Problem } from '@/lib/api'

export function ProblemsDrawer({ problems }: { problems: Problem[] }) {
  if (problems.length === 0) return null
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className="h-8 gap-1.5 border-amber-300 text-amber-800 hover:bg-amber-50"
        >
          <FileWarning className="size-3.5" />
          {problems.length} {problems.length === 1 ? 'problem' : 'problems'}
        </Button>
      </SheetTrigger>
      <SheetContent className="w-full overflow-y-auto data-[side=right]:sm:max-w-lg">
        <SheetHeader>
          <SheetTitle>Files fuda could not use</SheetTitle>
          <SheetDescription>
            These are not on the board. Fix them in the repo; the board picks them up on the next
            sync.
          </SheetDescription>
        </SheetHeader>
        <ul className="space-y-3 px-4 pb-4">
          {problems.map((p) => (
            <li
              key={p.path + p.reason}
              className="rounded-lg border border-amber-200 bg-amber-50/60 p-3"
            >
              <p className="font-mono text-xs break-all">{p.path}</p>
              <p className="mt-1 text-sm text-amber-900">{p.reason}</p>
            </li>
          ))}
        </ul>
      </SheetContent>
    </Sheet>
  )
}
