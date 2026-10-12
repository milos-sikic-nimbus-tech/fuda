// @vitest-environment jsdom
import { cleanup, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it } from 'vitest'
import { renderWithListing } from '@/test/render'
import type { HostStatus } from '@/lib/hosts'
import { AccountMenu } from './AccountMenu'

const github = (loggedIn: boolean): HostStatus => ({
  host: 'github',
  loggedIn,
  login: '/auth/github/login',
})
const azure = (loggedIn: boolean): HostStatus => ({
  host: 'azure',
  loggedIn,
  login: '/auth/azure/login',
})

async function openMenu(hosts: HostStatus[]) {
  renderWithListing(<AccountMenu />, { boards: [], hosts })
  await userEvent.click(screen.getByRole('button', { name: 'Account' }))
}

afterEach(cleanup)

describe('AccountMenu', () => {
  it('renders no button without hosts', () => {
    renderWithListing(<AccountMenu />, { boards: [], hosts: [] })
    expect(screen.queryByRole('button', { name: 'Account' })).toBeNull()
  })

  it('shows each host with its own state and action', async () => {
    await openMenu([github(true), azure(false)])
    const githubRow = screen.getByText('GitHub').closest<HTMLElement>('.flex')!
    const azureRow = screen.getByText('Azure DevOps').closest<HTMLElement>('.flex')!
    expect(within(githubRow).getByText('Logged in')).toBeTruthy()
    expect(within(githubRow).getByRole('button', { name: 'Log out' })).toBeTruthy()
    expect(within(azureRow).getByText('Logged out')).toBeTruthy()
    const logIn = within(azureRow).getByRole('link', { name: 'Log in' })
    expect(logIn.getAttribute('href')).toMatch(/^\/auth\/azure\/login\?return=/)
  })

  it.each([
    ['one host is logged in', [github(true), azure(false)], false],
    ['no host is logged in', [github(false), azure(false)], false],
    ['two hosts are logged in', [github(true), azure(true)], true],
  ])('Log out of all when %s: shown is %s', async (_name, hosts, shown) => {
    await openMenu(hosts)
    expect(screen.queryByText('Log out of all') !== null).toBe(shown)
  })
})
