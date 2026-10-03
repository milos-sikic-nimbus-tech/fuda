import { useNavigate, useSearch } from '@tanstack/react-router'
import { type BoardSearch, type ListKey, toggled } from '@/lib/filters'

export function useBoardSearch() {
  const search = useSearch({ from: '/' })
  const navigate = useNavigate({ from: '/' })

  const update = (patch: Partial<BoardSearch>) => {
    void navigate({ search: (prev) => ({ ...prev, ...patch }), replace: true })
  }

  const toggle = (key: ListKey, value: string) => update({ [key]: toggled(search[key], value) })

  const openTask = (id: string | undefined) => {
    void navigate({ search: (prev) => ({ ...prev, task: id }) })
  }

  const clear = () => {
    void navigate({
      search: (prev) => ({ task: prev.task, sort: prev.sort, view: prev.view }),
      replace: true,
    })
  }

  return { search, update, toggle, openTask, clear }
}
