<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import type { SupportedLocale } from '@/i18n'

const { locale } = useI18n()
const appStore = useAppStore()

function toggleLocale() {
  const newLocale: SupportedLocale = locale.value === 'zh-CN' ? 'en' : 'zh-CN'
  locale.value = newLocale
  appStore.setLocale(newLocale)
}

const displayLabel = () => {
  return locale.value === 'zh-CN' ? 'EN' : '中'
}
</script>

<template>
  <el-button
    size="small"
    round
    @click="toggleLocale"
    class="lang-switcher"
  >
    {{ displayLabel() }}
  </el-button>
</template>

<style scoped lang="scss">
.lang-switcher {
  font-weight: 600;
  min-width: 40px;
}
</style>
