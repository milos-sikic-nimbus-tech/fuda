import { RefreshCw, TriangleAlert } from 'lucide-react'
import { Hint } from '@/components/atoms/Hint'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import type { SyncStatus } from '@/lib/api'
import { useSync } from '@/lib/queries'
import { cn } from '@/lib/utils'

function ago(iso: string): string {
  const seconds = Math.max(0, Math.round((Date.now() - new Date(iso).getTime()) / 1000))
  if (seconds < 60) return 'just now'
  if (seconds < 3600) return `${Math.round(seconds / 60)} min ago`
  return `${Math.round(seconds / 3600)} h ago`
}

function shortSha(sha: string): string {
  return sha.startsWith('worktree-') ? 'working tree' : sha.slice(0, 8)
}

export function SyncIndicator({ sync }: { sync: SyncStatus }) {
  const mutation = useSync()
  const cooling = mutation.error?.message.includes('moments ago')

  return (
    <div className="flex items-center gap-2 text-xs text-muted-foreground">
      {sync.lastError && (
        <Tooltip>
          <TooltipTrigger asChild>
            <span className="inline-flex items-center gap-1 text-amber-700">
              <TriangleAlert className="size-3.5" />
              last sync failed
            </span>
          </TooltipTrigger>
          <TooltipContent className="max-w-80">{sync.lastError}</TooltipContent>
        </Tooltip>
      )}
      <Tooltip>
        <TooltipTrigger asChild>
          <span>
            {sync.develop.branch} @ <span className="font-mono">{shortSha(sync.develop.sha)}</span>{' '}
            · {ago(sync.develop.syncedAt)}
          </span>
        </TooltipTrigger>
        <TooltipContent>
          {sync.main
            ? sync.main.notYet
              ? `${sync.main.branch}: not yet (no board there)`
              : `${sync.main.branch} @ ${shortSha(sync.main.sha)}`
            : 'main is not watched'}
        </TooltipContent>
      </Tooltip>
      <Hint label={cooling ? 'Synced moments ago; try again shortly' : 'Read the repo now'}>
        <Button
          variant="ghost"
          size="icon"
          className="size-7"
          aria-label="Sync now"
          disabled={mutation.isPending}
          onClick={() => mutation.mutate()}
        >
          <RefreshCw className={cn('size-3.5', mutation.isPending && 'animate-spin')} />
        </Button>
      </Hint>
    </div>
  )
}
