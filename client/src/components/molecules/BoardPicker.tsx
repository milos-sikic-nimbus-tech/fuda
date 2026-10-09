import { Check, ChevronsUpDown } from 'lucide-react'
import { HostIcon } from '@/components/atoms/HostIcon'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { BOARD_PREFIX } from '@/lib/boardPath'
import { useBoards } from '@/lib/queries'

export function BoardPicker() {
  const { data: boards } = useBoards()
  if (!boards || boards.length < 2) return null
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="size-7" aria-label="Switch Board">
          <ChevronsUpDown className="size-3.5" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" className="min-w-48">
        {boards.map((board) => (
          <DropdownMenuItem key={board.path} asChild className="gap-2">
            <a href={board.path}>
              <HostIcon host={board.host} className="size-3.5 shrink-0" />
              <span className="truncate">{board.title}</span>
              {board.path === BOARD_PREFIX && <Check className="ml-auto size-3.5" />}
            </a>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
