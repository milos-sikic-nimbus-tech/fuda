import { Fragment } from 'react'
import { Check, ChevronsUpDown, ExternalLink } from 'lucide-react'
import { HostIcon } from '@/components/atoms/HostIcon'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import type { BoardData } from '@/lib/api'
import { BOARD_PREFIX } from '@/lib/boardPath'
import { groupBoards, pickerButtonLabel } from '@/lib/hosts'
import { useBoardListing } from '@/lib/queries'

export function WorkspacePicker({ board }: { board?: BoardData }) {
  const { data } = useBoardListing()
  if (!data || (data.boards.length === 0 && data.hosts.length === 0)) return null
  return (
    <>
      <span className="hidden text-border sm:inline">/</span>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="outline"
            className="h-7 min-w-0 max-w-40 gap-1.5 rounded-full px-2.5 text-[13px] shadow-xs md:max-w-none"
            aria-label="Switch Board"
          >
            {board && <HostIcon host={board.origin.host} className="size-3.5 shrink-0" />}
            <span className="truncate font-medium">{pickerButtonLabel(board)}</span>
            {board && (
              <span className="hidden font-normal text-muted-foreground md:inline">
                {board.sync.develop.branch}
              </span>
            )}
            <ChevronsUpDown className="size-3.5 shrink-0 text-muted-foreground" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start" className="min-w-48">
          {groupBoards(data.boards).map((group, index) => (
            <Fragment key={group.host}>
              {index > 0 && <DropdownMenuSeparator />}
              <DropdownMenuLabel>{group.label}</DropdownMenuLabel>
              {group.boards.map((entry) => {
                const current = entry.path === BOARD_PREFIX
                return (
                  <Fragment key={entry.path}>
                    <DropdownMenuItem asChild className="gap-2">
                      <a href={entry.path}>
                        <HostIcon host={entry.host} className="size-3.5 shrink-0" />
                        <span className="truncate">{entry.title}</span>
                        {current && <Check className="ml-auto size-3.5" />}
                      </a>
                    </DropdownMenuItem>
                    {current && board?.origin.url && (
                      <DropdownMenuItem asChild className="gap-2 pl-7 text-muted-foreground">
                        <a href={board.origin.url} target="_blank" rel="noreferrer">
                          <ExternalLink className="size-3.5" />
                          Open repo
                        </a>
                      </DropdownMenuItem>
                    )}
                  </Fragment>
                )
              })}
            </Fragment>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  )
}
