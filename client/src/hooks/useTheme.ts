import { useCallback, useEffect, useSyncExternalStore } from 'react'

export type ThemeChoice = 'light' | 'dark' | 'system'

const storageKey = 'fuda-theme'
const darkQuery = '(prefers-color-scheme: dark)'
const listeners = new Set<() => void>()

function readChoice(): ThemeChoice {
  try {
    const stored = localStorage.getItem(storageKey)
    return stored === 'light' || stored === 'dark' ? stored : 'system'
  } catch {
    return 'system'
  }
}

function resolve(choice: ThemeChoice): 'light' | 'dark' {
  if (choice !== 'system') return choice
  return window.matchMedia(darkQuery).matches ? 'dark' : 'light'
}

function apply(choice: ThemeChoice) {
  document.documentElement.classList.toggle('dark', resolve(choice) === 'dark')
}

function subscribe(listener: () => void) {
  listeners.add(listener)
  return () => listeners.delete(listener)
}

export function useTheme() {
  const choice = useSyncExternalStore(subscribe, readChoice, () => 'system' as const)

  useEffect(() => {
    apply(choice)
    if (choice !== 'system') return
    const media = window.matchMedia(darkQuery)
    const onChange = () => apply('system')
    media.addEventListener('change', onChange)
    return () => media.removeEventListener('change', onChange)
  }, [choice])

  const setChoice = useCallback((next: ThemeChoice) => {
    try {
      if (next === 'system') localStorage.removeItem(storageKey)
      else localStorage.setItem(storageKey, next)
    } catch {
      return
    }
    listeners.forEach((listener) => listener())
  }, [])

  return { choice, resolved: resolve(choice), setChoice }
}
