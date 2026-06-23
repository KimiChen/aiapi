const FULL_APP_ENTRY_PATH = '/app.html'
const DEFAULT_APP_PATH = '/dashboard'

export function currentBrowserPath(): string {
  return `${window.location.pathname}${window.location.search}${window.location.hash}`
}

export function normalizeFullAppTargetPath(value: string | null | undefined): string {
  if (!value) return DEFAULT_APP_PATH
  const target = value.trim()
  if (
    !target ||
    !target.startsWith('/') ||
    target.startsWith('//') ||
    target.includes('://') ||
    target.includes('\n') ||
    target.includes('\r')
  ) {
    return DEFAULT_APP_PATH
  }
  return target
}

export function fullAppEntryURL(targetPath = currentBrowserPath()): string {
  const url = new URL(FULL_APP_ENTRY_PATH, window.location.origin)
  url.searchParams.set('redirect', normalizeFullAppTargetPath(targetPath))
  return `${url.pathname}${url.search}`
}

export function navigateToFullApp(targetPath = currentBrowserPath(), replace = false): void {
  const next = fullAppEntryURL(targetPath)
  if (replace) {
    window.location.replace(next)
    return
  }
  window.location.assign(next)
}

export function consumeFullAppEntryRedirect(): void {
  const url = new URL(window.location.href)
  if (url.pathname !== FULL_APP_ENTRY_PATH && url.pathname !== `/static/app${FULL_APP_ENTRY_PATH}`) {
    return
  }
  const targetPath = normalizeFullAppTargetPath(url.searchParams.get('redirect'))
  window.history.replaceState(null, '', targetPath)
}
