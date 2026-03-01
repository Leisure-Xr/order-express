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
  background: var(--app-gradient-hero);
  border-radius: 18px;
  padding: 20px 18px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 20px rgba(231, 76, 60, 0.08);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    right: -30px;
    top: -30px;
    width: 120px;
    height: 120px;
    border-radius: 50%;
    background: rgba(231, 76, 60, 0.08);
    pointer-events: none;
  }

  &::after {
    content: '';
    position: absolute;
    right: 20px;
    bottom: -20px;
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: rgba(52, 152, 219, 0.06);
    pointer-events: none;
  }
}

.hero-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  position: relative;
  z-index: 1;
}

.name {
  font-size: 22px;
  font-weight: 900;
  background: linear-gradient(135deg, #303133, #e74c3c);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
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

.desc,
.addr {
  position: relative;
  z-index: 1;
}

.hero-actions {
  margin-top: 16px;
  display: flex;
  gap: 10px;
  align-items: center;
  position: relative;
  z-index: 1;
}

.section-title {
  font-weight: 900;
  font-size: 17px;
  color: #1a1a1a;
  margin-bottom: 10px;
  position: relative;
  padding-left: 12px;

  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 3px;
    bottom: 3px;
    width: 4px;
    border-radius: 2px;
    background: var(--app-gradient-primary);
  }
}

.rec-list {
  display: flex;
  flex-direction: column;
  gap: 12px;

  > :nth-child(1) { animation: fadeSlideUp 0.35s cubic-bezier(0.4, 0, 0.2, 1) 0.05s both; }
  > :nth-child(2) { animation: fadeSlideUp 0.35s cubic-bezier(0.4, 0, 0.2, 1) 0.10s both; }
  > :nth-child(3) { animation: fadeSlideUp 0.35s cubic-bezier(0.4, 0, 0.2, 1) 0.15s both; }
  > :nth-child(4) { animation: fadeSlideUp 0.35s cubic-bezier(0.4, 0, 0.2, 1) 0.20s both; }
}
</style>

