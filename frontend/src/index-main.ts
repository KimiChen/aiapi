import type { App as VueApp } from 'vue'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PublicApp from './PublicApp.vue'
import publicRouter, {
  isFullAppPublicPath,
  isGuestPublicPath,
} from './router/guest'
import i18n, { initPublicI18n } from './i18n'
import { useAppStore } from '@/stores/app'
import { enterFullApp, setPublicApp } from '@/public/fullAppBridge'
import './style.css'

function initThemeClass() {
  const savedTheme = localStorage.getItem('theme')
  const shouldUseDark =
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  document.documentElement.classList.toggle('dark', shouldUseDark)
}

function currentPath(): string {
  return `${window.location.pathname}${window.location.search}${window.location.hash}`
}

function hasPersistedAuthToken(): boolean {
  try {
    return Boolean(localStorage.getItem('auth_token'))
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
  const pinia = createPinia()
  app.use(pinia)

  const appStore = useAppStore()
  appStore.initFromInjectedConfig()

  if (appStore.siteName) {
    document.title = `${appStore.siteName} - Secure Portal`
  }

  await initPublicI18n()

  app.use(publicRouter)
  app.use(i18n)
  setPublicApp(app)

  await publicRouter.isReady()
  app.mount('#app')
  return app
}

async function bootstrap() {
  const pathname = window.location.pathname

  if (!isGuestPublicPath(pathname)) {
    if (hasPersistedAuthToken() || isFullAppPublicPath(pathname)) {
      await enterFullApp(currentPath())
      return
    }

    redirectToLoginWithReturnPath()
  }

  await mountPublicApp()
}

bootstrap()
