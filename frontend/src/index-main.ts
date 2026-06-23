import type { App as VueApp } from 'vue'
import { createApp } from 'vue'
import PublicApp from './PublicApp.vue'
import publicRouter, {
  isFullAppPublicPath,
  isGuestPublicPath,
} from './router/guest'
import { currentBrowserPath, navigateToFullApp } from '@/public/appNavigation'
import './style.css'

const GUEST_SITE_NAME = '企业数据中台'

function initThemeClass() {
  const savedTheme = localStorage.getItem('theme')
  const shouldUseDark =
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  document.documentElement.classList.toggle('dark', shouldUseDark)
}

function currentPath(): string {
  return currentBrowserPath()
}

function hasPersistedAuthSession(): boolean {
  try {
    const token = localStorage.getItem('auth_token')
    const rawUser = localStorage.getItem('auth_user')
    if (!token || !rawUser) {
      return false
    }
    const user = JSON.parse(rawUser)
    return Boolean(user && typeof user === 'object')
  } catch {
    return false
  }
}

function redirectToLoginWithReturnPath(): void {
  const redirect = currentPath()
  const url = new URL('/login', window.location.origin)
  if (redirect && redirect !== '/login') {
    url.searchParams.set('redirect', redirect)
  }
  window.history.replaceState(null, '', `${url.pathname}${url.search}${url.hash}`)
}

async function mountPublicApp(): Promise<VueApp<Element>> {
  initThemeClass()

  const app = createApp(PublicApp)

  document.title = `${GUEST_SITE_NAME} - Secure Portal`

  app.use(publicRouter)

  await publicRouter.isReady()
  app.mount('#app')
  return app
}

async function bootstrap() {
  const pathname = window.location.pathname

  if (!isGuestPublicPath(pathname)) {
    if (hasPersistedAuthSession() || isFullAppPublicPath(pathname)) {
      navigateToFullApp(currentPath(), true)
      return
    }

    redirectToLoginWithReturnPath()
  }

  await mountPublicApp()
}

bootstrap()
