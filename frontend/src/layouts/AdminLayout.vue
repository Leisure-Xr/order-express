<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { useLocaleText } from '@/composables/useLocaleText'
import LangSwitcher from '@/components/common/LangSwitcher.vue'
import {
  DataAnalysis,
  Food,
  List,
  Grid,
  Setting,
  Fold,
  Expand,
} from '@element-plus/icons-vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const settingsStore = useSettingsStore()
const { localText } = useLocaleText()

const storeName = computed(() => {
  const name = settingsStore.storeInfo ? localText(settingsStore.storeInfo.name) : ''
  return name || 'Order Express'
})

const storeLogo = computed(() => settingsStore.storeInfo?.logo || '')

function storeInitials(name: string): string {
  const trimmed = name.trim()
  if (!trimmed) return 'OE'
  const parts = trimmed.split(/\s+/).filter(Boolean)
  if (parts.length >= 2) {
    return `${parts[0]?.[0] ?? ''}${parts[1]?.[0] ?? ''}`.toUpperCase()
  }
  return trimmed.slice(0, 2).toUpperCase()
}

const storeMiniText = computed(() => storeInitials(storeName.value))

onMounted(() => {
  if (!settingsStore.storeInfo) {
    settingsStore.fetchStoreInfo().catch(() => {})
  }
})

function handleCommand(command: string) {
  if (command === 'logout') {
    authStore.logout()
    router.push({ name: 'AdminLogin' })
  }
}
</script>

<template>
  <el-container class="admin-layout">
    <el-aside :width="appStore.sidebarCollapsed ? '64px' : '220px'" class="admin-sidebar">
      <div class="logo-container">
        <div class="logo-inner">
          <el-avatar v-if="storeLogo" :size="32" :src="storeLogo" class="logo-avatar" />
          <span v-if="!appStore.sidebarCollapsed" class="logo-text">{{ storeName }}</span>
          <span v-else-if="!storeLogo" class="logo-text-mini">{{ storeMiniText }}</span>
        </div>
      </div>
      <el-menu
        :default-active="route.path"
        :collapse="appStore.sidebarCollapsed"
        :router="true"
        background-color="transparent"
        text-color="#cbd5e1"
        active-text-color="#e74c3c"
      >
        <el-menu-item index="/admin/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <template #title>{{ t('routes.dashboard') }}</template>
        </el-menu-item>
        <el-sub-menu index="menu-group">
          <template #title>
            <el-icon><Food /></el-icon>
            <span>{{ t('routes.menuManagement') }}</span>
          </template>
          <el-menu-item index="/admin/menu">{{ t('menu.allCategories') }}</el-menu-item>
          <el-menu-item index="/admin/menu/categories">{{ t('routes.categoryManagement') }}</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/admin/orders">
          <el-icon><List /></el-icon>
          <template #title>{{ t('routes.orderManagement') }}</template>
        </el-menu-item>
        <el-menu-item index="/admin/tables">
          <el-icon><Grid /></el-icon>
          <template #title>{{ t('routes.tableManagement') }}</template>
        </el-menu-item>
        <el-menu-item index="/admin/settings">
          <el-icon><Setting /></el-icon>
          <template #title>{{ t('routes.storeSettings') }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="admin-header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="appStore.toggleSidebar" :size="20">
            <Fold v-if="!appStore.sidebarCollapsed" />
            <Expand v-else />
          </el-icon>
        </div>
        <div class="header-right">
          <LangSwitcher />
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :src="authStore.user?.avatar">
                {{ authStore.user?.name?.charAt(0) }}
              </el-avatar>
              <span class="username">{{ authStore.user?.name }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">{{ t('auth.logout') }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="admin-main">
        <slot />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped lang="scss">
.admin-layout {
  height: 100vh;
}

.admin-sidebar {
  background: linear-gradient(180deg, #0b1220 0%, #111827 40%, #1f2937 100%);
  transition: width 0.3s ease;
  overflow: hidden;
  position: relative;

  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 200px;
    height: 200px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(231, 76, 60, 0.08) 0%, transparent 70%);
    pointer-events: none;
  }

  .logo-container {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  }

  .logo-inner {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    padding: 0 10px;
    min-width: 0;
  }

  .logo-avatar {
    border: 1px solid rgba(255, 255, 255, 0.22);
    box-shadow: 0 6px 18px rgba(0, 0, 0, 0.14);
    flex: none;
    background: rgba(255, 255, 255, 0.10);
  }

  .logo-text {
    font-size: 20px;
    font-weight: 800;
    white-space: nowrap;
    letter-spacing: 1px;
    overflow: hidden;
    text-overflow: ellipsis;
    background: linear-gradient(135deg, #ffffff, #e74c3c);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .logo-text-mini {
    font-size: 20px;
    font-weight: 800;
    background: linear-gradient(135deg, #ffffff, #e74c3c);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  :deep(.el-menu) {
    border-right: none;
  }

  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    position: relative;
    border-radius: 10px;
    margin: 4px 8px;
  }

  :deep(.el-menu-item:hover),
  :deep(.el-sub-menu__title:hover) {
    background: rgba(255, 255, 255, 0.06) !important;
  }

  :deep(.el-menu-item.is-active) {
    background: rgba(231, 76, 60, 0.14) !important;
  }

  :deep(.el-menu-item.is-active::after) {
    content: '';
    position: absolute;
    left: 0;
    top: 8px;
    bottom: 8px;
    width: 3px;
    border-radius: 999px;
    background: linear-gradient(180deg, #e74c3c, #f0574a);
  }

  :deep(.el-sub-menu.is-active > .el-sub-menu__title) {
    background: rgba(231, 76, 60, 0.10) !important;
  }

  :deep(.el-menu--collapse) {
    width: 64px;
  }
}

.admin-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 60px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(2, 6, 23, 0.06);
}

.header-left {
  display: flex;
  align-items: center;

  .collapse-btn {
    cursor: pointer;
    color: #606266;
    padding: 6px;
    border-radius: 8px;
    transition: color 0.2s, background 0.2s;

    &:hover {
      color: #e74c3c;
      background: rgba(231, 76, 60, 0.06);
    }
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  outline: none;

  .username {
    font-size: 14px;
    color: #333333;
  }
}

.admin-main {
  padding: 24px;
  background:
    radial-gradient(1200px 700px at 15% 0%, rgba(231, 76, 60, 0.10), transparent 55%),
    radial-gradient(1000px 600px at 100% 0%, rgba(52, 152, 219, 0.08), transparent 50%),
    var(--el-bg-color-page);
  overflow-y: auto;
}
</style>
