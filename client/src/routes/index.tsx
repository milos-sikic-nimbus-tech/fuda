import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: BoardPage,
})

function BoardPage() {
  return <main className="p-6 text-sm text-muted-foreground">fuda</main>
}
