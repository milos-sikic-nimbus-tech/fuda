export type Column = {
  id: string
  name: string
  statuses: string[]
  prOpen: boolean
  unknown: boolean
}

export type Card = {
  id: string
  title: string
  status: string
  column: string
  owners: string[]
  testers: string[]
  labels: string[]
  prefix: string
  added: string
  claimed: string
  done?: string
  prs: number[]
  prRef?: string
  openPrs: number[]
  blockedBy?: string
  blockedByIds: string[]
  blocks: string[]
  references: string[]
  referencedBy: string[]
  inProd: boolean
  newInPr: boolean
  archiveCandidate: boolean
  path: string
}

export type LabelGroup = { name: string; values: string[] }

export type Facets = {
  statuses: string[]
  people: string[]
  labelGroups: LabelGroup[]
  prefixes: string[]
}

export type Problem = { path: string; reason: string }

export type BranchStatus = { branch: string; sha: string; syncedAt: string; notYet?: boolean }

export type SyncStatus = {
  develop: BranchStatus
  main?: BranchStatus
  lastAttempt: string
  lastError?: string
}

export type BoardData = {
  title: string
  prLink?: string
  columns: Column[]
  cards: Card[]
  facets: Facets
  problems: Problem[]
  sync: SyncStatus
}

export type TaskDetail = Card & {
  custom: { key: string; value: string }[]
  html: string
}

export type DocDetail = { path: string; title: string; html: string; backlinks: string[] }

export class ApiError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

async function request<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, init)
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new ApiError(res.status, body.error ?? res.statusText)
  }
  return res.status === 202 ? (undefined as T) : res.json()
}

export const api = {
  board: () => request<BoardData>('/api/board'),
  task: (id: string) => request<TaskDetail>(`/api/tasks/${encodeURIComponent(id)}`),
  search: (q: string) => request<string[]>(`/api/search?q=${encodeURIComponent(q)}`),
  archive: () => request<Card[]>('/api/archive'),
  doc: (path: string) => request<DocDetail>(`/api/docs?path=${encodeURIComponent(path)}`),
  sync: () => request<void>('/api/sync', { method: 'POST' }),
}

export function prHref(template: string | undefined, n: number): string | undefined {
  return template?.replace('{n}', String(n))
}
