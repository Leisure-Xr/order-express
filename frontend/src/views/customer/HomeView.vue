<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useSettingsStore } from '@/stores/settings'
import { useMenuStore } from '@/stores/menu'
import { useAppStore } from '@/stores/app'
import { useLocaleText } from '@/composables/useLocaleText'
import DishCard from '@/components/customer/DishCard.vue'

const { t } = useI18n()
const router = useRouter()
const settingsStore = useSettingsStore()
const menuStore = useMenuStore()
const appStore = useAppStore()
const { localText } = useLocaleText()

const loading = ref(false)

const storeName = computed(() => (settingsStore.storeInfo ? localText(settingsStore.storeInfo.name) : 'Order Express'))
const storeDesc = computed(() => (settingsStore.storeInfo ? localText(settingsStore.storeInfo.description) : ''))
const storeAddr = computed(() => (settingsStore.storeInfo ? localText(settingsStore.storeInfo.address) : ''))

const isOpen = computed(() => settingsStore.isStoreOpen())

const recommended = computed(() => {
  const list = menuStore.dishes.filter((d) => d.status === 'on_sale')
  const preferred = list.filter((d) => d.tags?.includes('popular') || d.tags?.includes('new'))
  return (preferred.length ? preferred : list).slice(0, 4)
})

onMounted(async () => {
  loading.value = true
  try {
    if (!settingsStore.storeInfo) {
      await settingsStore.fetchStoreInfo()
      await settingsStore.fetchBusinessHours()
    }
    if (!menuStore.dishes.length) {
      await menuStore.fetchAll()
    }
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    loading.value = false
  }
})

function startOrder() {
  router.push({ name: 'CustomerMenu' })
}

function openDish(dish: any) {
  router.push({ name: 'CustomerDishDetail', params: { id: dish.id } })
}
</script>

<template>
  <div class="home-page">
    <el-skeleton v-if="loading && !settingsStore.storeInfo" :rows="6" animated />

    <div v-else class="content">
      <div class="hero">
        <div class="hero-top">
          <div class="name">{{ storeName }}</div>
          <el-tag :type="isOpen ? 'success' : 'danger'">
            {{ isOpen ? 'OPEN' : 'CLOSED' }}
          </el-tag>
        </div>
        <div v-if="storeDesc" class="desc">{{ storeDesc }}</div>
        <div v-if="storeAddr" class="addr">{{ storeAddr }}</div>

        <div class="hero-actions">
          <el-button type="primary" size="large" @click="startOrder">{{ t('routes.menu') }}</el-button>
          <el-tag v-if="appStore.currentTableId" type="warning" effect="dark">
            {{ t('order.tableNumber') }}: {{ appStore.currentTableId }}
          </el-tag>
        </div>
      </div>

      <div class="section">
        <div class="section-title">{{ t('dashboard.popularItems') }}</div>
        <div class="rec-list">
          <DishCard
            v-for="dish in recommended"
            :key="dish.id"
            :dish="dish"
            :show-add="false"
            @click="openDish"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.home-page {
  padding: 12px;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.hero {
  background: linear-gradient(135deg, rgba(231, 76, 60, 0.16), rgba(52, 152, 219, 0.10));
  border-radius: 16px;
  padding: 16px;
  border: 1px solid rgba(0, 0, 0, 0.04);
}

.hero-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.name {
  font-size: 20px;
  font-weight: 900;
  color: #303133;
}

.desc {
  margin-top: 10px;
  color: #606266;
  font-size: 13px;
  line-height: 1.6;
}

.addr {
  margin-top: 8px;
  color: #909399;
  font-size: 12px;
}

.hero-actions {
  margin-top: 12px;
  display: flex;
  gap: 10px;
  align-items: center;
}

.section-title {
  font-weight: 900;
  color: #303133;
  margin-bottom: 10px;
}

.rec-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>

