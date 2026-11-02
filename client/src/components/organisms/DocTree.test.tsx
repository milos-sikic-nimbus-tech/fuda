// @vitest-environment jsdom
import {
  createMemoryHistory,
  createRootRoute,
  createRoute,
  createRouter,
  RouterProvider,
} from '@tanstack/react-router'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it } from 'vitest'
import { DocTree } from './DocTree'

afterEach(cleanup)

async function renderTree(openPath: string) {
  const root = createRootRoute({
    component: () => (
      <DocTree paths={['README.md', 'board/stages.md', 'guides/deep/one.md']} openPath={openPath} />
    ),
  })
  const splat = createRoute({ getParentRoute: () => root, path: '/docs/$' })
  const router = createRouter({
    routeTree: root.addChildren([splat]),
    history: createMemoryHistory({ initialEntries: ['/docs/' + openPath] }),
  })
  render(<RouterProvider router={router} />)
  await screen.findByRole('navigation', { name: 'Docs' })
}

describe('DocTree', () => {
  it('highlights the open doc and expands its folders', async () => {
    await renderTree('guides/deep/one.md')
    expect(screen.getByRole('link', { name: 'one.md' }).getAttribute('aria-current')).toBe('page')
    expect(screen.getByRole('link', { name: 'README.md' }).getAttribute('aria-current')).toBeNull()
  })

  it('collapses and expands a folder', async () => {
    await renderTree('README.md')
    expect(screen.getByRole('link', { name: 'stages.md' })).toBeTruthy()
    await userEvent.click(screen.getByRole('button', { name: 'board' }))
    expect(screen.queryByRole('link', { name: 'stages.md' })).toBeNull()
    await userEvent.click(screen.getByRole('button', { name: 'board' }))
    expect(screen.getByRole('link', { name: 'stages.md' })).toBeTruthy()
  })
})
