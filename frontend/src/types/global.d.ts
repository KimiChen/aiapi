import type { PublicSettingsConfig } from '@/types'

declare global {
  interface Window {
    __APP_CONFIG__?: PublicSettingsConfig
  }
}

export {}
