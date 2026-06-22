<template>
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <div v-else class="min-h-screen bg-wiki-bg font-body text-wiki-txt">
    <div
      v-if="mobileOpen"
      class="fixed inset-0 z-40 bg-black/30 md:hidden"
      @click="mobileOpen = false"
    ></div>

    <div class="flex min-h-screen">
      <aside
        class="fixed inset-y-0 left-0 z-50 flex w-60 flex-col border-r border-wiki-border bg-white transition-transform duration-200 md:sticky md:top-0 md:translate-x-0"
        :class="mobileOpen ? 'translate-x-0' : '-translate-x-full'"
      >
        <div class="border-b border-wiki-border p-5">
          <RouterLink to="/home" class="flex items-center gap-3">
            <div
              class="flex h-9 w-9 items-center justify-center overflow-hidden rounded-xl bg-gradient-to-br from-wiki-accent to-wiki-accent2"
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
              <h1 class="truncate font-heading text-base font-bold text-wiki-txt">
                {{ displayName }}
              </h1>
              <p class="text-[11px] uppercase tracking-wide text-wiki-muted">DATA FABRIC</p>
            </div>
          </RouterLink>
        </div>

        <nav class="flex-1 py-3">
          <a
            v-for="item in navItems"
            :key="item.label"
            class="guest-nav-item relative flex cursor-pointer items-center gap-3 px-5 py-2.5 text-sm"
            :class="{ active: item.active }"
            @click="handleNavClick(item.target)"
          >
            <Icon :name="item.icon" size="sm" :stroke-width="2" />
            <span>{{ item.label }}</span>
          </a>
        </nav>

        <div class="border-t border-wiki-border p-4">
          <button
            type="button"
            class="mb-3 flex items-center gap-2 text-xs text-wiki-muted transition-colors hover:text-wiki-txt"
          >
            <Icon name="globe" size="xs" :stroke-width="2" />
            <span>English</span>
          </button>
          <RouterLink
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="flex items-center gap-2 text-xs text-wiki-muted transition-colors hover:text-wiki-txt"
          >
            <Icon name="arrowLeft" size="xs" :stroke-width="2" />
            <span>{{ isAuthenticated ? '进入控制台' : '登录客户区' }}</span>
          </RouterLink>
        </div>
      </aside>

      <main class="min-w-0 flex-1">
        <header
          class="sticky top-0 z-30 flex items-center justify-between border-b border-wiki-border bg-white/80 px-4 py-3 backdrop-blur-lg md:px-6"
        >
          <div class="flex min-w-0 items-center gap-3">
            <button
              type="button"
              class="rounded-lg p-1.5 text-wiki-muted transition-colors hover:bg-wiki-surface2 md:hidden"
              @click="mobileOpen = true"
              aria-label="Open menu"
            >
              <Icon name="menu" size="md" :stroke-width="2" />
            </button>
            <h2 class="font-heading text-lg font-semibold text-wiki-txt">数据中台</h2>
          </div>

          <div class="flex items-center gap-3">
            <div class="hidden items-center gap-2 rounded-lg bg-wiki-surface2 px-3 py-1.5 text-sm sm:flex">
              <Icon name="shield" size="sm" class="text-wiki-accent" />
              <span class="font-semibold">统一认证</span>
              <span class="text-wiki-muted">已启用</span>
            </div>
            <RouterLink
              :to="isAuthenticated ? dashboardPath : '/login'"
              class="inline-flex items-center justify-center gap-2 rounded-xl bg-wiki-accent px-4 py-2.5 text-sm font-semibold text-white transition-all hover:-translate-y-0.5 hover:bg-indigo-600 hover:shadow-[0_4px_12px_rgba(99,102,241,0.3)]"
            >
              <Icon :name="isAuthenticated ? 'grid' : 'login'" size="sm" :stroke-width="2" />
              <span>{{ isAuthenticated ? '控制台' : '登录' }}</span>
            </RouterLink>
          </div>
        </header>

        <div class="mx-auto max-w-7xl px-4 py-6 md:px-6">
          <section class="grid gap-4 md:grid-cols-3">
            <div
              v-for="card in overviewCards"
              :key="card.label"
              class="rounded-xl border border-wiki-border bg-white p-5 transition-all hover:-translate-y-0.5 hover:shadow-[0_4px_16px_rgba(0,0,0,0.06)]"
            >
              <div
                class="mb-4 flex h-10 w-10 items-center justify-center rounded-lg"
                :class="card.iconClass"
              >
                <Icon :name="card.icon" size="sm" :stroke-width="2" />
              </div>
              <p class="text-sm font-medium text-wiki-muted">{{ card.label }}</p>
              <p class="mt-1 font-heading text-2xl font-semibold text-wiki-txt">{{ card.value }}</p>
              <p class="mt-1 text-xs text-wiki-muted">{{ card.detail }}</p>
            </div>
          </section>

          <section class="mt-5 grid gap-5 lg:grid-cols-[minmax(0,1.4fr)_minmax(320px,0.8fr)]">
            <div class="rounded-xl border border-wiki-border bg-white p-5">
              <div class="mb-5 flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
                <div>
                  <h3 class="font-heading text-base font-semibold text-wiki-txt">
                    数据资产概览（近 30 天）
                  </h3>
                  <p class="mt-1 text-xs text-wiki-muted">今日同步 24 批次 · 异常 0 项</p>
                </div>
                <RouterLink
                  to="/login"
                  class="inline-flex items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-indigo-500 to-violet-600 px-4 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-200 transition-all hover:from-indigo-600 hover:to-violet-700"
                >
                  <Icon name="login" size="sm" :stroke-width="2" />
                  <span>申请接入</span>
                </RouterLink>
              </div>

              <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-5">
                <div
                  v-for="metric in usageMetrics"
                  :key="metric.label"
                  class="rounded-lg border border-slate-100 bg-slate-50/50 p-4"
                >
                  <div class="mb-3 flex items-center gap-2">
                    <Icon :name="metric.icon" size="xs" :class="metric.color" :stroke-width="2" />
                    <span class="text-xs font-medium text-wiki-muted">{{ metric.label }}</span>
                  </div>
                  <p class="font-heading text-xl font-semibold text-wiki-txt">{{ metric.value }}</p>
                  <p class="mt-1 text-[11px] text-wiki-muted">{{ metric.detail }}</p>
                </div>
              </div>
            </div>

            <div
              class="rounded-2xl border border-indigo-200/80 bg-gradient-to-r from-indigo-500 to-violet-600 p-5 text-white shadow-lg shadow-indigo-500/15"
            >
              <div class="mb-5 flex items-center justify-between gap-4">
                <div>
                  <h3 class="font-heading text-base font-semibold">平台运行态</h3>
                  <p class="mt-1 text-xs text-white/65">实时概览</p>
                </div>
                <div class="rounded-full bg-white/15 p-2">
                  <Icon name="chart" size="sm" :stroke-width="2" />
                </div>
              </div>

              <div class="grid gap-3">
                <div
                  v-for="metric in platformMetrics"
                  :key="metric.label"
                  class="rounded-xl bg-white/10 p-4 ring-1 ring-white/10"
                >
                  <div class="flex items-center justify-between gap-3">
                    <span class="text-sm text-white/75">{{ metric.label }}</span>
                    <Icon :name="metric.icon" size="xs" class="text-white/70" :stroke-width="2" />
                  </div>
                  <p class="mt-2 font-heading text-2xl font-semibold">{{ metric.value }}</p>
                  <p class="mt-1 text-xs text-white/60">{{ metric.detail }}</p>
                </div>
              </div>
            </div>
          </section>

          <section class="mt-5 rounded-xl border border-wiki-border bg-white p-5">
            <div class="mb-4 flex flex-col justify-between gap-3 sm:flex-row sm:items-center">
              <div>
                <h3 class="font-heading text-base font-semibold text-wiki-txt">数据服务目录</h3>
                <p class="text-xs text-wiki-muted">
                  按主题域查看来源系统、同步策略与服务等级
                </p>
              </div>
              <div class="rounded-lg bg-wiki-surface2 p-1 text-xs font-medium text-wiki-muted">
                <span class="rounded-md bg-white px-3 py-1 text-wiki-txt shadow-sm">24小时</span>
                <span class="px-3 py-1">7天</span>
              </div>
            </div>

            <div class="overflow-x-auto rounded-lg border border-wiki-border">
              <table class="w-full min-w-[720px] text-left text-sm">
                <thead class="bg-wiki-surface2 text-xs uppercase tracking-wide text-wiki-muted">
                  <tr>
                    <th class="px-4 py-3 font-semibold">主题域</th>
                    <th class="px-4 py-3 font-semibold">数据源</th>
                    <th class="px-4 py-3 font-semibold">数据量</th>
                    <th class="px-4 py-3 font-semibold">同步策略</th>
                    <th class="px-4 py-3 font-semibold">最近校验</th>
                    <th class="px-4 py-3 text-right font-semibold">服务等级</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-wiki-border">
                  <tr
                    v-for="row in modelRows"
                    :key="row.domain"
                    class="transition-colors hover:bg-wiki-bg"
                  >
                    <td class="px-4 py-3 font-medium text-wiki-txt">{{ row.domain }}</td>
                    <td class="px-4 py-3 text-wiki-muted">{{ row.sources }}</td>
                    <td class="px-4 py-3 text-wiki-muted">{{ row.output }}</td>
                    <td class="px-4 py-3 text-wiki-muted">{{ row.work }}</td>
                    <td class="px-4 py-3 text-wiki-muted">{{ row.latency }}</td>
                    <td class="px-4 py-3 text-right font-semibold text-wiki-accent">{{ row.sla }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>

          <section class="mt-5 grid gap-5 lg:grid-cols-[minmax(0,1fr)_360px]">
            <div class="rounded-xl border border-wiki-border bg-white p-5">
              <div class="mb-4 flex items-center gap-2">
                <Icon name="book" size="sm" class="text-wiki-accent" :stroke-width="2" />
                <h3 class="font-heading text-base font-semibold text-wiki-txt">接入规范</h3>
              </div>
              <p class="text-sm text-wiki-muted">
                面向业务系统提供统一数据模型、逻辑数据仓库与实时 CDC 管道，降低跨库查询与报表整合成本。
              </p>
              <div class="mt-4 rounded-xl bg-slate-900 p-4 font-mono text-xs text-slate-100">
                <p class="text-emerald-300">SELECT customer_id, risk_level, updated_at</p>
                <p class="pl-3 text-slate-300">FROM fabric.customer_profile</p>
                <p class="pl-3 text-slate-300">WHERE biz_date = CURRENT_DATE</p>
                <p class="pl-3 text-amber-200">AND data_quality_status = 'PASS';</p>
              </div>
            </div>

            <div class="rounded-xl border border-wiki-border bg-white p-5">
              <div class="mb-4 flex items-center justify-between gap-3">
                <div>
                  <h3 class="font-heading text-base font-semibold text-wiki-txt">最近任务</h3>
                  <p class="text-xs text-wiki-muted">查看全部</p>
                </div>
                <Icon name="clock" size="sm" class="text-slate-300" :stroke-width="1.5" />
              </div>
              <div class="flex min-h-[130px] flex-col items-center justify-center text-center">
                <Icon name="clock" size="xl" class="mb-3 text-slate-200" :stroke-width="1.5" />
                <p class="text-sm text-wiki-muted">暂无异常任务</p>
              </div>
            </div>
          </section>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore, useAppStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

const mobileOpen = ref(false)

const siteName = computed(() => {
  const name = appStore.cachedPublicSettings?.site_name || appStore.siteName || ''
  const trimmed = name.trim()
  const legacyDefaultName = ['Sub', '2', 'API'].join('')
  return trimmed && trimmed !== legacyDefaultName ? trimmed : ''
})
const displayName = computed(() => siteName.value || '数据中台')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))

const navItems = [
  { label: '总览', target: 'overview', icon: 'grid' as const, active: true },
  { label: '数据目录', target: 'catalog', icon: 'database' as const, active: false },
  { label: '数据治理', target: 'governance', icon: 'shield' as const, active: false },
  { label: '交换任务', target: 'usage', icon: 'sync' as const, active: false },
  { label: '服务编排', target: 'models', icon: 'cube' as const, active: false },
  { label: '接入规范', target: 'docs', icon: 'book' as const, active: false }
]

const overviewCards = [
  {
    label: '已登记数据资产',
    value: '1,248 项',
    detail: '覆盖 18 个业务域',
    icon: 'database' as const,
    iconClass: 'bg-indigo-50 text-wiki-accent'
  },
  {
    label: '数据服务',
    value: '86 个服务',
    detail: '统一目录与权限',
    icon: 'cube' as const,
    iconClass: 'bg-violet-50 text-wiki-accent2'
  },
  {
    label: '质量规则',
    value: '312 条在线',
    detail: '自动稽核与追溯',
    icon: 'shield' as const,
    iconClass: 'bg-purple-50 text-purple-500'
  }
]

const usageMetrics = [
  { label: '交换批次', value: '724', detail: '今日已完成', icon: 'sync' as const, color: 'text-indigo-500' },
  { label: '实时通道', value: '36', detail: 'CDC 运行中', icon: 'download' as const, color: 'text-sky-500' },
  { label: '服务调用', value: '18.6K', detail: '业务系统访问', icon: 'upload' as const, color: 'text-emerald-500' },
  { label: '质量稽核', value: '99.6%', detail: '规则通过率', icon: 'database' as const, color: 'text-amber-500' },
  { label: '数据血缘', value: '1,904', detail: '字段级链路', icon: 'calculator' as const, color: 'text-violet-600' }
]

const platformMetrics = [
  { label: '数据新鲜度', value: '96.8%', detail: '按 SLA 达标', icon: 'database' as const },
  { label: '平均同步延迟', value: '2.4s', detail: '实时链路', icon: 'clock' as const },
  { label: '可用性', value: '99.95%', detail: '核心服务', icon: 'bolt' as const }
]

const modelRows = [
  { domain: '客户主数据', sources: 'CRM / 会员中心', output: '18.2M 行', work: 'CDC 实时', latency: '2 分钟前', sla: '核心' },
  { domain: '交易明细', sources: '订单库 / 支付库', output: '42.6M 行', work: '准实时', latency: '1 分钟前', sla: '核心' },
  { domain: '库存与供应链', sources: 'ERP / WMS', output: '7.4M 行', work: '15 分钟批', latency: '通过', sla: '重要' },
  { domain: '经营指标', sources: '财务 / BI', output: '864 张宽表', work: '小时级', latency: '通过', sla: '重要' },
  { domain: '风控标签', sources: '规则引擎 / 日志', output: '2,318 标签', work: '事件驱动', latency: '通过', sla: '敏感' }
]

function handleNavClick(target: string) {
  mobileOpen.value = false
  if (target === 'login') {
    router.push('/login')
  }
}

onMounted(() => {
  authStore.checkAuth()
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
