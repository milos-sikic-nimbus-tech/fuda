import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/guide/$')({
  component: GuidePage,
})

function GuidePage() {
  return (
    <main className="p-6 text-sm text-muted-foreground">
      The guide is coming in the next milestone.
    </main>
  )
}
