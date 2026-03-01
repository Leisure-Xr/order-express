<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useSettingsStore } from '@/stores/settings'
import { useLocaleText } from '@/composables/useLocaleText'
import { getOrderHistoryIds } from '@/utils/orderHistory'
import { getOrderByIdApi } from '@/api/order'
import {
  Document,
  Lock,
  UserFilled,
  ArrowRight,
} from '@element-plus/icons-vue'

const { t } = useI18n()
const router = useRouter()
const appStore = useAppStore()
const settingsStore = useSettingsStore()
const { localText } = useLocaleText()

const orderCount = ref(0)
const totalSpent = ref(0)
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const ids = getOrderHistoryIds()
    orderCount.value = ids.length

    if (!settingsStore.storeInfo) {
      await settingsStore.fetchStoreInfo()
    }

    if (ids.length > 0) {
      const results = await Promise.all(
        ids.slice(0, 20).map(async (id) => {
          try {
            const res = await getOrderByIdApi(id)
            return res.code === 200 && res.data ? res.data.total : 0
          } catch {
            return 0
          }
        }),
      )
      totalSpent.value = results.reduce((sum, v) => sum + v, 0)
    }
  } catch {
    // Silently fail — profile is informational
  } finally {
    loading.value = false
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

function goOrders() {
  router.push({ name: 'CustomerOrderHistory' })
}
</script>

<template>
  <div class="profile-page">
    <div class="title">{{ t('profile.title') }}</div>

    <el-skeleton v-if="loading" :rows="6" animated />

    <template v-else>
      <el-card shadow="never" class="section-card">
        <div class="summary-row">
          <div class="summary-item">
            <div class="summary-value">{{ orderCount }}</div>
            <div class="summary-label">{{ t('profile.totalOrders') }}</div>
          </div>
          <div class="summary-divider" />
          <div class="summary-item">
            <div class="summary-value">&yen;{{ totalSpent.toFixed(0) }}</div>
            <div class="summary-label">{{ t('profile.totalSpent') }}</div>
          </div>
        </div>
      </el-card>

      <el-card shadow="never" class="section-card menu-card">
        <div class="menu-item" @click="goOrders">
          <div class="menu-left">
            <div class="menu-icon" style="background: rgba(231, 76, 60, 0.10); color: #e74c3c">
              <el-icon :size="18"><Document /></el-icon>
            </div>
            <span class="menu-label">{{ t('profile.myOrders') }}</span>
          </div>
          <el-icon class="menu-arrow" :size="16"><ArrowRight /></el-icon>
        </div>

        <div class="menu-divider" />

        <div class="menu-item">
          <div class="menu-left">
            <div class="menu-icon" style="background: rgba(52, 152, 219, 0.10); color: #3498db">
              <el-icon :size="18"><Lock /></el-icon>
            </div>
            <span class="menu-label">{{ t('profile.privacy') }}</span>
          </div>
          <el-icon class="menu-arrow" :size="16"><ArrowRight /></el-icon>
        </div>

        <div class="menu-divider" />

        <div class="menu-item">
          <div class="menu-left">
            <div class="menu-icon" style="background: rgba(243, 156, 18, 0.10); color: #f39c12">
              <el-icon :size="18"><UserFilled /></el-icon>
            </div>
            <span class="menu-label">{{ t('profile.account') }}</span>
          </div>
          <el-icon class="menu-arrow" :size="16"><ArrowRight /></el-icon>
        </div>
      </el-card>

      <el-card v-if="appStore.currentTableId" shadow="never" class="section-card">
        <div class="setting-row">
          <span class="setting-label">{{ t('order.tableNumber') }}</span>
          <div class="setting-actions">
            <el-tag type="warning" effect="dark">{{ appStore.currentTableId }}</el-tag>
            <el-button size="small" text type="danger" @click="clearTable">
              {{ t('profile.clearTable') }}
            </el-button>
          </div>
        </div>
      </el-card>

      <el-card shadow="never" class="section-card">
        <div class="setting-row">
          <span class="setting-label">{{ t('profile.language') }}</span>
          <div class="setting-actions">
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
      </el-card>

      <el-card shadow="never" class="section-card">
        <div class="setting-row">
          <span class="setting-label">{{ t('profile.storeName') }}</span>
          <span class="setting-text">{{ storeName }}</span>
        </div>
      </el-card>
    </template>
  </div>
</template>

<style scoped lang="scss">
.profile-page {
  padding: 12px;
}

.title {
  font-size: 18px;
  font-weight: 900;
  color: #303133;
  margin-bottom: 12px;
}

.section-card {
  border-radius: 14px;
  border: 1px solid rgba(0, 0, 0, 0.03);
  margin-bottom: 12px;
}

.summary-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
}

.summary-item {
  flex: 1;
  text-align: center;
  padding: 8px 0;
}

.summary-value {
  font-size: 24px;
  font-weight: 900;
  background: linear-gradient(135deg, #e74c3c, #f0574a);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.summary-label {
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
}

.summary-divider {
  width: 1px;
  height: 40px;
  background: rgba(0, 0, 0, 0.06);
  flex: none;
}

.menu-card {
  :deep(.el-card__body) {
    padding: 4px 16px;
  }
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 0;
  cursor: pointer;
  transition: opacity var(--app-transition-fast);

  &:active {
    opacity: 0.6;
  }
}

.menu-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.menu-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex: none;
}

.menu-label {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.menu-arrow {
  color: #c0c4cc;
}

.menu-divider {
  height: 1px;
  background: rgba(0, 0, 0, 0.04);
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.setting-label {
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

.setting-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.setting-text {
  color: #909399;
  font-size: 14px;
}
</style>
