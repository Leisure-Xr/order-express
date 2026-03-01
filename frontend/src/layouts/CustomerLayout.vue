<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'
import { useCartStore } from '@/stores/cart'
import {
  HomeFilled,
  Food,
  ShoppingCart,
  User,
} from '@element-plus/icons-vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const cartStore = useCartStore()

const tabs = computed(() => [
  { path: '/', icon: HomeFilled, label: t('routes.home') },
  { path: '/menu', icon: Food, label: t('routes.menu') },
  { path: '/cart', icon: ShoppingCart, label: t('routes.cart'), badge: cartStore.itemCount },
  { path: '/profile', icon: User, label: t('routes.profile') },
])

function isActive(tabPath: string): boolean {
  if (tabPath === '/') {
    return route.path === '/'
  }
  return route.path.startsWith(tabPath)
}

function navigateTo(path: string) {
  router.push(path)
}

const tableDisplay = computed(() => {
  if (appStore.currentTableId) {
    return appStore.currentTableId
  }
  return null
})
</script>

<template>
  <div class="customer-layout">
    <header class="customer-header">
      <div class="customer-header-inner">
        <div class="header-store-name">Order Express</div>
        <div class="header-actions">
          <el-tag v-if="tableDisplay" type="warning" size="small" effect="dark">
            {{ t('order.tableNumber') }}: {{ tableDisplay }}
          </el-tag>
        </div>
      </div>
    </header>

    <main class="customer-content">
      <div class="customer-content-inner">
        <slot />
      </div>
    </main>

    <nav class="customer-bottom-bar">
      <div class="customer-bottom-inner">
        <div
          v-for="tab in tabs"
          :key="tab.path"
          class="tab-item"
          :class="{ active: isActive(tab.path) }"
          @click="navigateTo(tab.path)"
        >
          <div class="tab-icon-wrapper">
            <el-badge
              v-if="tab.badge && tab.badge > 0"
              :value="tab.badge"
              :max="99"
              class="tab-badge"
            >
              <el-icon :size="22"><component :is="tab.icon" /></el-icon>
            </el-badge>
            <el-icon v-else :size="22"><component :is="tab.icon" /></el-icon>
          </div>
          <span class="tab-label">{{ tab.label }}</span>
        </div>
      </div>
    </nav>
  </div>
</template>

<style scoped lang="scss">
.customer-layout {
  --customer-tabbar-height: 56px;
  --customer-safe-bottom: env(safe-area-inset-bottom, 0px);
  --customer-bottom-bar-height: calc(var(--customer-tabbar-height) + var(--customer-safe-bottom));

  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: var(--el-bg-color-page);
}

.customer-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--app-glass-bg);
  backdrop-filter: var(--app-glass-blur);
  border-bottom: 1px solid rgba(2, 6, 23, 0.06);
}

.customer-header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 0 auto;
  width: 100%;
  max-width: var(--app-content-max-width-sm);
  padding: 0 12px;
  height: 48px;
}

.header-store-name {
  font-size: 18px;
  font-weight: 800;
  background: linear-gradient(135deg, #e74c3c, #f0574a);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 0.5px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.customer-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: calc(var(--customer-bottom-bar-height) + 12px);
}

.customer-content-inner {
  margin: 0 auto;
  width: 100%;
  max-width: var(--app-content-max-width-sm);
}

.customer-bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  background: var(--app-glass-bg);
  backdrop-filter: var(--app-glass-blur);
  border-top: 1px solid rgba(2, 6, 23, 0.06);
  height: var(--customer-bottom-bar-height);
  padding-bottom: var(--customer-safe-bottom);
}

.customer-bottom-inner {
  margin: 0 auto;
  width: 100%;
  max-width: var(--app-content-max-width-sm);
  height: var(--customer-tabbar-height);
  display: flex;
  align-items: stretch;
  justify-content: space-around;
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  height: 100%;
  cursor: pointer;
  color: #999999;
  transition: color 0.2s;
  user-select: none;
  -webkit-tap-highlight-color: transparent;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    inset: 6px 10px;
    border-radius: 14px;
    background: rgba(231, 76, 60, 0.10);
    opacity: 0;
    transition: opacity 0.2s ease;
  }

  &.active {
    color: #e74c3c;

    &::before {
      opacity: 1;
    }

    .tab-icon-wrapper {
      transform: scale(1.1);
    }
  }

  &:active {
    opacity: 0.7;
  }
}

.tab-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 24px;
  position: relative;
  z-index: 1;
  transition: transform var(--app-transition-base);
}

.tab-label {
  font-size: 11px;
  margin-top: 2px;
  line-height: 1;
  position: relative;
  z-index: 1;
}

.tab-badge {
  :deep(.el-badge__content) {
    font-size: 10px;
    height: 16px;
    line-height: 16px;
    padding: 0 4px;
  }
}
</style>
