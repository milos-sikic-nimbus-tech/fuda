import { Link } from '@tanstack/react-router'

export function TaskLink({ id }: { id: string }) {
  return (
    <Link
      to="/"
      search={(prev) => ({ ...prev, task: id })}
      className="font-mono text-xs text-primary underline-offset-2 hover:underline"
    >
      {id}
    </Link>
  )
}
