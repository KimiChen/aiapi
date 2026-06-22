<template>
  <div class="min-h-screen bg-halo-bg font-body text-halo-txt">
    <div class="flex min-h-screen flex-col lg:flex-row">
      <aside
        class="hidden w-[260px] flex-col border-r border-halo-border bg-white lg:flex"
      >
        <div class="border-b border-halo-border p-5">
          <RouterLink to="/home" class="flex items-center gap-3">
            <div
              class="flex h-9 w-9 items-center justify-center overflow-hidden rounded-xl bg-gradient-to-br from-halo-accent to-halo-accent2"
            >
              <img
                v-if="siteLogo"
                :src="siteLogo"
                alt=""
                class="h-full w-full object-contain"
              />
              <Icon v-else name="sparkles" size="md" class="text-white" :stroke-width="2" />
            </div>
            <div class="min-w-0">
              <h1 class="truncate font-heading text-base font-bold text-halo-txt">
                {{ displayName }}
              </h1>
              <p class="text-[11px] uppercase tracking-wide text-halo-muted">API SERVICE</p>
            </div>
          </RouterLink>
        </div>

        <nav class="flex-1 py-3">
          <RouterLink
            v-for="item in navItems"
            :key="item.label"
            :to="item.to"
            class="guest-nav-item relative flex items-center gap-3 px-5 py-2.5 text-sm"
            :class="{ active: item.active }"
          >
            <Icon :name="item.icon" size="sm" :stroke-width="2" />
            <span>{{ item.label }}</span>
          </RouterLink>
        </nav>

        <div class="border-t border-halo-border p-4">
          <RouterLink
            to="/home"
            class="flex items-center gap-2 text-xs text-halo-muted transition-colors hover:text-halo-txt"
          >
            <Icon name="arrowLeft" size="xs" :stroke-width="2" />
            <span>返回门户</span>
          </RouterLink>
        </div>
      </aside>

      <main class="flex min-h-screen flex-1 flex-col">
        <header
          class="sticky top-0 z-20 flex items-center justify-between border-b border-halo-border bg-white/80 px-4 py-3 backdrop-blur-lg sm:px-6"
        >
          <RouterLink to="/home" class="flex min-w-0 items-center gap-3 lg:hidden">
            <div
              class="flex h-9 w-9 items-center justify-center overflow-hidden rounded-xl bg-gradient-to-br from-halo-accent to-halo-accent2"
            >
              <img
                v-if="siteLogo"
                :src="siteLogo"
                alt=""
                class="h-full w-full object-contain"
              />
              <Icon v-else name="sparkles" size="md" class="text-white" :stroke-width="2" />
            </div>
            <div class="min-w-0">
              <p class="truncate font-heading text-base font-semibold text-halo-txt">
                {{ displayName }}
              </p>
              <p class="text-[11px] uppercase tracking-wide text-halo-muted">API SERVICE</p>
            </div>
          </RouterLink>

          <div class="hidden items-center gap-2 rounded-lg bg-halo-surface2 px-3 py-1.5 text-sm lg:flex">
            <Icon name="dollar" size="sm" class="text-halo-accent" />
            <span class="font-semibold">0.00</span>
            <span class="text-halo-muted">credits</span>
          </div>

          <div class="flex items-center gap-2">
            <a
              v-if="docUrl"
              :href="docUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex h-9 w-9 items-center justify-center rounded-lg text-halo-muted transition-colors hover:bg-halo-surface2 hover:text-halo-txt"
              aria-label="Docs"
            >
              <Icon name="book" size="sm" :stroke-width="2" />
            </a>
            <RouterLink
              to="/home"
              class="hidden rounded-lg border border-halo-border px-3 py-1.5 text-xs font-semibold text-halo-muted transition-colors hover:bg-halo-surface2 hover:text-halo-txt sm:inline-flex"
            >
              概览
            </RouterLink>
          </div>
        </header>

        <section class="flex flex-1 items-center justify-center px-4 py-8 sm:px-6 lg:px-10">
          <div class="grid w-full max-w-6xl gap-6 lg:grid-cols-[minmax(0,1fr)_430px]">
            <div class="order-2 min-w-0 space-y-5 lg:order-1">
              <div
                class="rounded-2xl border border-indigo-200/80 bg-gradient-to-r from-indigo-500 to-violet-600 p-6 text-white shadow-lg shadow-indigo-500/15"
              >
                <div class="mb-5 flex items-center justify-between gap-4">
                  <div>
                    <p class="text-sm font-medium text-white/75">Unified API Service</p>
                    <h2 class="mt-1 font-heading text-2xl font-semibold">Gateway Console</h2>
                  </div>
                  <div class="rounded-full bg-white/15 p-2">
                    <Icon name="sparkles" size="md" :stroke-width="2" />
                  </div>
                </div>
                <div class="grid gap-3 sm:grid-cols-3">
                  <div
                    v-for="metric in heroMetrics"
                    :key="metric.label"
                    class="rounded-xl bg-white/10 p-4 ring-1 ring-white/10"
                  >
                    <p class="text-[11px] uppercase tracking-wide text-white/60">{{ metric.label }}</p>
                    <p class="mt-2 font-heading text-2xl font-semibold">{{ metric.value }}</p>
                    <p class="mt-1 text-xs text-white/65">{{ metric.detail }}</p>
                  </div>
                </div>
              </div>

              <div class="grid gap-4 md:grid-cols-3">
                <div
                  v-for="card in statCards"
                  :key="card.label"
                  class="rounded-xl border border-halo-border bg-white p-4 transition-all hover:-translate-y-0.5 hover:shadow-[0_4px_16px_rgba(0,0,0,0.06)]"
                >
                  <div
                    class="mb-4 flex h-9 w-9 items-center justify-center rounded-lg"
                    :class="card.iconClass"
                  >
                    <Icon :name="card.icon" size="sm" :stroke-width="2" />
                  </div>
                  <p class="text-xs font-medium text-halo-muted">{{ card.label }}</p>
                  <p class="mt-1 font-heading text-xl font-semibold text-halo-txt">{{ card.value }}</p>
                  <p class="mt-1 text-xs text-halo-muted">{{ card.detail }}</p>
                </div>
              </div>

              <div class="rounded-xl border border-halo-border bg-white p-5">
                <div class="mb-4 flex items-center justify-between gap-3">
                  <div>
                    <h3 class="font-heading text-base font-semibold text-halo-txt">Model TPS</h3>
                    <p class="text-xs text-halo-muted">近 24 小时平台指标</p>
                  </div>
                  <div class="rounded-lg bg-halo-surface2 p-1 text-xs font-medium text-halo-muted">
                    <span class="rounded-md bg-white px-3 py-1 text-halo-txt shadow-sm">24h</span>
                    <span class="px-3 py-1">7d</span>
                  </div>
                </div>

                <div class="overflow-x-auto rounded-lg border border-halo-border">
                  <table class="w-full min-w-[560px] text-left text-sm">
                    <thead class="bg-halo-surface2 text-xs uppercase tracking-wide text-halo-muted">
                      <tr>
                        <th class="px-4 py-3 font-semibold">Model</th>
                        <th class="px-4 py-3 font-semibold">Requests</th>
                        <th class="px-4 py-3 font-semibold">Output</th>
                        <th class="px-4 py-3 font-semibold">Latency</th>
                        <th class="px-4 py-3 text-right font-semibold">TPS</th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-halo-border">
                      <tr
                        v-for="row in modelRows"
                        :key="row.model"
                        class="transition-colors hover:bg-halo-bg"
                      >
                        <td class="px-4 py-3 font-medium text-halo-txt">{{ row.model }}</td>
                        <td class="px-4 py-3 text-halo-muted">{{ row.requests }}</td>
                        <td class="px-4 py-3 text-halo-muted">{{ row.output }}</td>
                        <td class="px-4 py-3 text-halo-muted">{{ row.latency }}</td>
                        <td class="px-4 py-3 text-right font-semibold text-halo-accent">{{ row.tps }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            <div class="order-1 min-w-0 lg:order-2">
              <div class="rounded-2xl border border-halo-border bg-white p-6 shadow-[0_8px_24px_rgba(15,23,42,0.05)] sm:p-8">
                <slot />
              </div>

              <div class="mt-5 text-center text-sm">
                <slot name="footer" />
              </div>

              <p class="mt-6 text-center text-xs text-slate-400">
                &copy; {{ currentYear }} {{ displayName }}. All rights reserved.
              </p>
            </div>
          </div>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useAppStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'

const appStore = useAppStore()

const siteName = computed(() => {
  const name = appStore.cachedPublicSettings?.site_name || appStore.siteName || ''
  return name.trim() && name.trim() !== 'Sub2API' ? name.trim() : ''
})
const displayName = computed(() => siteName.value || 'AI Console')
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const currentYear = computed(() => new Date().getFullYear())

const navItems = [
  { label: '概览', to: '/home', icon: 'grid' as const, active: false },
  { label: 'API Key', to: '/login', icon: 'key' as const, active: true },
  { label: '充值', to: '/login', icon: 'creditCard' as const, active: false },
  { label: '用量', to: '/login', icon: 'chart' as const, active: false },
  { label: '模型', to: '/login', icon: 'cpu' as const, active: false },
  { label: '文档', to: '/home', icon: 'book' as const, active: false }
]

const heroMetrics = [
  { label: 'Cache Hit', value: '88.1%', detail: 'read + write' },
  { label: 'Latency', value: '11.1s', detail: 'per request' },
  { label: 'Requests', value: '6,425', detail: 'last 24h' }
]

const statCards = [
  {
    label: 'Token Balance',
    value: '0.000000',
    detail: 'credits available',
    icon: 'dollar' as const,
    iconClass: 'bg-indigo-50 text-halo-accent'
  },
  {
    label: 'API Key',
    value: '1 active',
    detail: 'auto routing enabled',
    icon: 'key' as const,
    iconClass: 'bg-violet-50 text-halo-accent2'
  },
  {
    label: 'Models',
    value: '20 online',
    detail: 'multi-platform access',
    icon: 'cpu' as const,
    iconClass: 'bg-purple-50 text-purple-500'
  }
]

const modelRows = [
  { model: 'claude-opus-4', requests: '2,164', output: '19.32M', latency: '11.3s', tps: '791.8' },
  { model: 'claude-sonnet-4', requests: '644', output: '36.59M', latency: '11.7s', tps: '1267.3' },
  { model: 'gpt-5.5', requests: '1,839', output: '584.1K', latency: '10.6s', tps: '29.9' },
  { model: 'gpt-image-2', requests: '173', output: '438', latency: '38.3s', tps: '19.2' }
]

onMounted(() => {
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.guest-nav-item {
  color: #64748b;
  transition: all 0.2s ease;
}

.guest-nav-item:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.guest-nav-item.active {
  background: rgba(99, 102, 241, 0.08);
  color: #6366f1;
  font-weight: 600;
}

.guest-nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 8px;
  bottom: 8px;
  width: 3px;
  border-radius: 0 3px 3px 0;
  background: #6366f1;
}
</style>
