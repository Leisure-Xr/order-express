<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { useLocaleText } from '@/composables/useLocaleText'
import {
  Box,
  CreditCard,
  Document,
  Grid,
  Lock,
  Service,
  Shop,
  Switch,
  SwitchButton,
  Van,
  ArrowRight,
} from '@element-plus/icons-vue'

const { t } = useI18n()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const settingsStore = useSettingsStore()
const { localText } = useLocaleText()

onMounted(async () => {
  try {
    if (!settingsStore.storeInfo) {
      await settingsStore.fetchStoreInfo()
    }
  } catch {
    // Silently fail — profile is informational
  }
})

const storeName = computed(() =>
  settingsStore.storeInfo ? localText(settingsStore.storeInfo.name) : 'Order Express',
)

function switchLocale(val: string | number) {
  const newLocale = val as 'zh-CN' | 'en'
  appStore.setLocale(newLocale)
}

function clearTable() {
  appStore.setTableFromQuery(null)
}

type OrderQuickFilter = 'unpaid' | 'to_ship' | 'to_receive' | 'after_sales'

function goOrders(filter?: 'all' | OrderQuickFilter) {
  router.push({
    name: 'CustomerOrderHistory',
    query: filter && filter !== 'all' ? { filter } : undefined,
  })
}

function goUserAgreement() {
  router.push({ name: 'CustomerUserAgreement' })
}

function goPrivacyPolicy() {
  router.push({ name: 'CustomerPrivacyPolicy' })
}

async function confirmLogout() {
  try {
    await ElMessageBox.confirm(t('profile.logoutConfirm'), t('profile.logout'), {
      type: 'warning',
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
    })
  } catch {
    return
  }

  await authStore.logout()
  ElMessage.success(t('common.success'))
  router.push({ name: 'CustomerHome' })
}
</script>

<template>
  <div class="profile-page">
    <el-card shadow="never" class="section-card orders-card">
      <div class="section-header">
        <div class="section-title">{{ t('profile.myOrders') }}</div>
        <div class="section-link" @click="goOrders('all')">
          <span>{{ t('profile.allOrders') }}</span>
          <el-icon :size="16"><ArrowRight /></el-icon>
        </div>
      </div>

      <div class="order-shortcuts">
        <div class="shortcut" @click="goOrders('unpaid')">
          <el-icon class="shortcut-icon" :size="26"><CreditCard /></el-icon>
          <div class="shortcut-label">{{ t('profile.toPay') }}</div>
        </div>
        <div class="shortcut" @click="goOrders('to_ship')">
          <el-icon class="shortcut-icon" :size="26"><Box /></el-icon>
          <div class="shortcut-label">{{ t('profile.toShip') }}</div>
        </div>
        <div class="shortcut" @click="goOrders('to_receive')">
          <el-icon class="shortcut-icon" :size="26"><Van /></el-icon>
          <div class="shortcut-label">{{ t('profile.toReceive') }}</div>
        </div>
        <div class="shortcut" @click="goOrders('after_sales')">
          <el-icon class="shortcut-icon" :size="26"><Service /></el-icon>
          <div class="shortcut-label">{{ t('profile.afterSales') }}</div>
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="section-card list-card">
      <div class="section-title">{{ t('profile.privacy') }}</div>

      <div class="cell cell-clickable" @click="goUserAgreement">
        <div class="cell-left">
          <el-icon class="cell-icon" :size="20"><Document /></el-icon>
          <span class="cell-label">{{ t('profile.userAgreement') }}</span>
        </div>
        <el-icon class="cell-arrow" :size="16"><ArrowRight /></el-icon>
      </div>

      <div class="cell-divider" />

      <div class="cell cell-clickable" @click="goPrivacyPolicy">
        <div class="cell-left">
          <el-icon class="cell-icon" :size="20"><Lock /></el-icon>
          <span class="cell-label">{{ t('profile.privacyPolicy') }}</span>
        </div>
        <el-icon class="cell-arrow" :size="16"><ArrowRight /></el-icon>
      </div>
    </el-card>

    <el-card shadow="never" class="section-card list-card">
      <div class="section-title">{{ t('profile.account') }}</div>

      <div class="cell">
        <div class="cell-left">
          <el-icon class="cell-icon" :size="20"><Switch /></el-icon>
          <span class="cell-label">{{ t('profile.language') }}</span>
        </div>
        <div class="cell-right">
          <el-segmented
            :model-value="appStore.locale"
            :options="[
              { label: '中文', value: 'zh-CN' },
              { label: 'English', value: 'en' },
            ]"
            @change="switchLocale"
          />
        </div>
      </div>

      <template v-if="appStore.currentTableId">
        <div class="cell-divider" />
        <div class="cell">
          <div class="cell-left">
            <el-icon class="cell-icon" :size="20"><Grid /></el-icon>
            <span class="cell-label">{{ t('order.tableNumber') }}</span>
          </div>
          <div class="cell-right">
            <el-tag type="warning" effect="dark">{{ appStore.currentTableId }}</el-tag>
            <el-button size="small" text type="danger" @click="clearTable">
              {{ t('profile.clearTable') }}
            </el-button>
          </div>
        </div>
      </template>

      <div class="cell-divider" />

      <div class="cell">
        <div class="cell-left">
          <el-icon class="cell-icon" :size="20"><Shop /></el-icon>
          <span class="cell-label">{{ t('profile.storeName') }}</span>
        </div>
        <div class="cell-right">
          <span class="cell-value">{{ storeName }}</span>
        </div>
      </div>

      <template v-if="authStore.isAuthenticated">
        <div class="cell-divider" />
        <div class="cell cell-clickable logout-row" @click="confirmLogout">
          <div class="cell-left">
            <el-icon class="cell-icon" :size="20"><SwitchButton /></el-icon>
            <span class="cell-label">{{ t('profile.logout') }}</span>
          </div>
          <el-icon class="cell-arrow" :size="16"><ArrowRight /></el-icon>
        </div>
      </template>
    </el-card>
  </div>
</template>

<style scoped lang="scss">
.profile-page {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-card {
  border-radius: 14px;
  border: 1px solid rgba(0, 0, 0, 0.03);
  :deep(.el-card__body) {
    padding: 14px 16px;
  }
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 10px;
}

.section-title {
  font-size: 18px;
  font-weight: 900;
  color: #303133;
}

.section-link {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #c0c4cc;
  font-size: 14px;
  cursor: pointer;
  user-select: none;

  &:active {
    opacity: 0.6;
  }
}

.order-shortcuts {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 6px;
}

.shortcut {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 12px 0;
  border-radius: 12px;
  cursor: pointer;
  transition: opacity var(--app-transition-fast), background var(--app-transition-fast);

  &:active {
    opacity: 0.6;
  }
}

.shortcut-icon {
  color: #111827;
}

.shortcut-label {
  margin-top: 8px;
  font-size: 13px;
  color: #303133;
}

.list-card {
  .section-title {
    padding-bottom: 6px;
  }
}

.cell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 0;
}

.cell-clickable {
  cursor: pointer;
  user-select: none;
  transition: opacity var(--app-transition-fast);

  &:active {
    opacity: 0.6;
  }
}

.cell-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.cell-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cell-icon {
  color: #111827;
}

.cell-label {
  font-size: 15px;
  color: #303133;
  font-weight: 600;
}

.cell-value {
  color: #909399;
  font-size: 14px;
}

.cell-arrow {
  color: #c0c4cc;
}

.cell-divider {
  height: 1px;
  background: rgba(0, 0, 0, 0.06);
}

.logout-row {
  .cell-label {
    color: #f56c6c;
  }
}
</style>
