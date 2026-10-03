import { RefreshCw } from 'lucide-react'
import { useEffect, useRef } from 'react'
import { toast } from 'sonner'
import { Hint } from '@/components/atoms/Hint'
import { Button } from '@/components/ui/button'
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
  return sha.startsWith('worktree-') ? 'working tree' : sha.slice(0, 7)
}

function useSyncErrorToast(lastError: string | undefined) {
  const shown = useRef<string | undefined>(undefined)
  useEffect(() => {
    if (lastError && lastError !== shown.current) {
      toast.warning('The last sync failed', {
        description: `${lastError}. Showing the last good copy.`,
      })
    }
    shown.current = lastError
  }, [lastError])
}

export function SyncIndicator({ sync }: { sync: SyncStatus }) {
  const mutation = useSync()
  useSyncErrorToast(sync.lastError)

  const mainNote = sync.main
    ? sync.main.notYet
      ? `${sync.main.branch}: no board there yet`
      : `${sync.main.branch} @ ${shortSha(sync.main.sha)}`
    : 'main is not watched'

  return (
    <div className="flex items-center gap-2 text-xs text-muted-foreground">
      <Hint
        label={
          <>
            Last read {ago(sync.develop.syncedAt)} · {mainNote}
          </>
        }
      >
        <span className="inline-flex items-center gap-2">
          <span
            className={cn(
              'size-1.5 rounded-full ring-3',
              sync.lastError ? 'bg-amber-500 ring-amber-100' : 'bg-primary ring-primary/15',
            )}
          />
          <span className="hidden lg:inline">synced {ago(sync.develop.syncedAt)}</span>
          <span className="hidden font-mono text-[11px] lg:inline">
            {shortSha(sync.develop.sha)}
          </span>
        </span>
      </Hint>
      <Hint label="Read the repo now">
        <Button
          variant="outline"
          size="sm"
          className="h-7 gap-1.5 bg-card px-2 text-xs shadow-xs"
          aria-label="Sync now"
          disabled={mutation.isPending}
          onClick={() => mutation.mutate()}
        >
          <RefreshCw className={cn('size-3.5', mutation.isPending && 'animate-spin')} />
          <span className="hidden sm:inline">Sync</span>
        </Button>
      </Hint>
    </div>
  )
}
