import { Check, ChevronsUpDown } from 'lucide-react'
import { HostIcon } from '@/components/atoms/HostIcon'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { BOARD_PREFIX, loginPath } from '@/lib/boardPath'
import { hostLabel, logOut } from '@/lib/hosts'
import { useBoardListing } from '@/lib/queries'

export function BoardPicker() {
  const { data } = useBoardListing()
  if (!data || (data.boards.length === 0 && data.hosts.length === 0)) return null
  const { boards, hosts } = data
  const returnTo = window.location.pathname + window.location.search
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
        {hosts.length > 0 && <DropdownMenuSeparator />}
        {hosts.map((host) =>
          host.loggedIn ? (
            <DropdownMenuItem key={host.host} onSelect={() => void logOut(host.host)}>
              Log out of {hostLabel(host.host)}
            </DropdownMenuItem>
          ) : (
            <DropdownMenuItem key={host.host} asChild>
              <a href={loginPath(returnTo, host.login)}>Log in to {hostLabel(host.host)}</a>
            </DropdownMenuItem>
          ),
        )}
        {hosts.length > 1 && hosts.some((host) => host.loggedIn) && (
          <DropdownMenuItem onSelect={() => void logOut()}>Log out of all</DropdownMenuItem>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
