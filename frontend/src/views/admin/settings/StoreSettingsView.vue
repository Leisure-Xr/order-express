<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useSettingsStore } from '@/stores/settings'
import type { StoreInfo, BusinessHours, DeliverySettings } from '@/types'

const { t } = useI18n()
const settingsStore = useSettingsStore()

const activeTab = ref<'basic' | 'hours' | 'delivery'>('basic')

const storeInfoModel = ref<StoreInfo | null>(null)
const businessHoursModel = ref<BusinessHours[]>([])
const deliveryModel = ref<DeliverySettings | null>(null)

function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

onMounted(async () => {
  try {
    await settingsStore.fetchAllSettings()
    storeInfoModel.value = settingsStore.storeInfo ? clone(settingsStore.storeInfo) : null
    businessHoursModel.value = clone(settingsStore.businessHours)
    deliveryModel.value = settingsStore.deliverySettings ? clone(settingsStore.deliverySettings) : null
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  }
})

function dayLabel(day: number): string {
  const labels = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
  return labels[day] || String(day)
}

async function saveBasic() {
  if (!storeInfoModel.value) return
  try {
    await settingsStore.updateStoreInfo(storeInfoModel.value)
    ElMessage.success(t('common.success'))
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  }
}

async function saveHours() {
  try {
    await settingsStore.updateBusinessHours(businessHoursModel.value)
    ElMessage.success(t('common.success'))
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  }
}

async function saveDelivery() {
  if (!deliveryModel.value) return
  try {
    await settingsStore.updateDeliverySettings(deliveryModel.value)
    ElMessage.success(t('common.success'))
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  }
}
</script>

<template>
  <div class="settings-page">
    <div class="title">{{ t('routes.storeSettings') }}</div>

    <el-tabs v-model="activeTab">
      <el-tab-pane :label="t('settings.basicInfo')" name="basic">
        <el-card shadow="never" class="card">
          <el-form v-if="storeInfoModel" label-position="top">
            <div class="grid">
              <el-form-item :label="t('settings.shopName') + ' (ZH)'" class="col-6">
                <el-input v-model="storeInfoModel.name.zh" />
              </el-form-item>
              <el-form-item :label="t('settings.shopName') + ' (EN)'" class="col-6">
                <el-input v-model="storeInfoModel.name.en" />
              </el-form-item>
            </div>

            <div class="grid">
              <el-form-item :label="t('settings.address') + ' (ZH)'" class="col-6">
                <el-input v-model="storeInfoModel.address.zh" />
              </el-form-item>
              <el-form-item :label="t('settings.address') + ' (EN)'" class="col-6">
                <el-input v-model="storeInfoModel.address.en" />
              </el-form-item>
            </div>

            <div class="grid">
              <el-form-item :label="t('settings.phone')" class="col-6">
                <el-input v-model="storeInfoModel.phone" />
              </el-form-item>
              <el-form-item :label="t('settings.shopLogo')" class="col-6">
                <el-input v-model="storeInfoModel.logo" placeholder="https://..." />
              </el-form-item>
            </div>

            <div class="grid">
              <el-form-item :label="t('menu.description') + ' (ZH)'" class="col-6">
                <el-input v-model="storeInfoModel.description.zh" type="textarea" :rows="3" />
              </el-form-item>
              <el-form-item :label="t('menu.description') + ' (EN)'" class="col-6">
                <el-input v-model="storeInfoModel.description.en" type="textarea" :rows="3" />
              </el-form-item>
            </div>

            <el-button type="primary" :loading="settingsStore.loading" @click="saveBasic">
              {{ t('common.save') }}
            </el-button>
          </el-form>
        </el-card>
      </el-tab-pane>

      <el-tab-pane :label="t('settings.businessHours')" name="hours">
        <el-card shadow="never" class="card">
          <el-table :data="businessHoursModel" stripe>
            <el-table-column label="Day" width="120">
              <template #default="{ row }">{{ dayLabel(row.dayOfWeek) }}</template>
            </el-table-column>
            <el-table-column :label="t('settings.closed')" width="120">
              <template #default="{ row }">
                <el-switch v-model="row.isClosed" />
              </template>
            </el-table-column>
            <el-table-column :label="t('settings.startTime')">
              <template #default="{ row }">
                <el-input v-model="row.openTime" placeholder="09:00" :disabled="row.isClosed" />
              </template>
            </el-table-column>
            <el-table-column :label="t('settings.endTime')">
              <template #default="{ row }">
                <el-input v-model="row.closeTime" placeholder="21:00" :disabled="row.isClosed" />
              </template>
            </el-table-column>
          </el-table>

          <div class="actions">
            <el-button type="primary" :loading="settingsStore.loading" @click="saveHours">
              {{ t('common.save') }}
            </el-button>
          </div>
        </el-card>
      </el-tab-pane>

      <el-tab-pane :label="t('settings.deliverySettings')" name="delivery">
        <el-card shadow="never" class="card">
          <el-form v-if="deliveryModel" label-position="top">
            <el-form-item :label="t('settings.enableDelivery')">
              <el-switch v-model="deliveryModel.enabled" />
            </el-form-item>

            <div class="grid">
              <el-form-item :label="t('settings.minOrderAmount')" class="col-4">
                <el-input-number v-model="deliveryModel.minimumOrder" :min="0" :step="1" />
              </el-form-item>
              <el-form-item :label="t('settings.deliveryFee')" class="col-4">
                <el-input-number v-model="deliveryModel.deliveryFee" :min="0" :step="1" />
              </el-form-item>
              <el-form-item :label="t('settings.freeDeliveryThreshold')" class="col-4">
                <el-input-number v-model="deliveryModel.freeDeliveryThreshold" :min="0" :step="1" />
              </el-form-item>
            </div>

            <div class="grid">
              <el-form-item :label="t('settings.estimatedDeliveryTime')" class="col-6">
                <el-input-number v-model="deliveryModel.estimatedMinutes" :min="0" :step="1" />
              </el-form-item>
              <el-form-item :label="t('settings.deliveryRange')" class="col-6">
                <el-input-number v-model="deliveryModel.deliveryRadius" :min="0" :step="1" />
              </el-form-item>
            </div>

            <el-button type="primary" :loading="settingsStore.loading" @click="saveDelivery">
              {{ t('common.save') }}
            </el-button>
          </el-form>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped lang="scss">
.settings-page {
  max-width: 1100px;
}

.title {
  font-size: 18px;
  font-weight: 800;
  color: #303133;
  margin-bottom: 12px;
}

.card {
  border-radius: 12px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 12px;
}

.col-4 {
  grid-column: span 4;
}

.col-6 {
  grid-column: span 6;
}

.actions {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 900px) {
  .grid {
    grid-template-columns: repeat(1, 1fr);
  }

  .col-4,
  .col-6 {
    grid-column: span 1;
  }
}
</style>
