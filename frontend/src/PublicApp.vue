<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import { watch } from 'vue'
import Toast from '@/components/common/Toast.vue'
import NavigationProgress from '@/components/common/NavigationProgress.vue'
import { useAppStore } from '@/stores/app'
import { resolveDocumentTitle } from '@/router/title'

const route = useRoute()
const appStore = useAppStore()

function updateDocumentTitle() {
  document.title = resolveDocumentTitle(route.meta.title, appStore.siteName, route.meta.titleKey as string)
}

watch(
  [
    () => route.fullPath,
    () => route.meta.title,
    () => route.meta.titleKey,
    () => appStore.siteName,
  ],
  updateDocumentTitle,
  { immediate: true }
)
</script>

<template>
  <NavigationProgress />
  <RouterView />
  <Toast />
</template>
