import type { BoardListing } from './api'

export type HostName = 'github' | 'azure'

export type HostStatus = { host: HostName; loggedIn: boolean; login: string; error?: string }

export type BoardListData = { boards: BoardListing[]; hosts: HostStatus[] }

const labels: Record<HostName, string> = { github: 'GitHub', azure: 'Azure DevOps' }

export function hostLabel(host: HostName): string {
  return labels[host]
}

export function loginRedirect({ hosts }: BoardListData): string | undefined {
  const [only] = hosts
  return hosts.length === 1 && !only?.loggedIn ? only?.login : undefined
}

export function onlyBoard({ boards, hosts }: BoardListData): BoardListing | undefined {
  return boards.length === 1 && hosts.every((host) => host.loggedIn) ? boards[0] : undefined
}

export async function logOut(host?: HostName) {
  await fetch(host ? `/auth/${host}/logout` : '/auth/logout', { method: 'POST' })
  window.location.assign('/')
}
