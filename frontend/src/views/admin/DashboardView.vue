<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, type RouteLocationRaw } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, Food, List, Grid, Setting } from '@element-plus/icons-vue'
import { useOrderStore } from '@/stores/order'
import { useTableStore } from '@/stores/table'
import { useLocaleText } from '@/composables/useLocaleText'
import AdminPageHeader from '@/components/admin/AdminPageHeader.vue'

const { t } = useI18n()
const router = useRouter()
const { localText } = useLocaleText()
const orderStore = useOrderStore()
const tableStore = useTableStore()

const loading = ref(false)

function formatDateParam(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const today = new Date()
const sevenDaysAgo = new Date()
sevenDaysAgo.setDate(today.getDate() - 6)

const dateRange = ref<[Date, Date]>([sevenDaysAgo, today])

const shortcuts = computed(() => [
  {
    text: t('dashboard.today'),
    value: () => {
      const d = new Date()
      return [d, d] as [Date, Date]
    },
  },
  {
    text: t('dashboard.last7Days'),
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setDate(start.getDate() - 6)
      return [start, end] as [Date, Date]
    },
  },
  {
    text: t('dashboard.last30Days'),
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setDate(start.getDate() - 29)
      return [start, end] as [Date, Date]
    },
  },
])

async function loadStats() {
  loading.value = true
  try {
    const [from, to] = dateRange.value
    await Promise.all([
      orderStore.fetchOrderStats({
        from: formatDateParam(from),
        to: formatDateParam(to),
      }),
      tableStore.fetchTables(),
    ])
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    loading.value = false
  }
}

onMounted(loadStats)
watch(dateRange, loadStats)

const stats = computed(() => orderStore.orderStats)
const totalOrders = computed(() => stats.value?.totalOrders ?? 0)
const paidRevenue = computed(() => stats.value?.paidRevenue ?? 0)
const pendingOrders = computed(() => stats.value?.pendingOrders ?? 0)
const tablesInUse = computed(() => tableStore.occupiedTables.length)

const revenueChart = computed(() => {
  const days = (stats.value?.dailyRevenue ?? []).map((d) => ({
    date: d.date.slice(5),
    amount: d.amount,
  }))
  const max = Math.max(...days.map((d) => d.amount), 1)
  return { days, max }
})

const popularItems = computed(() => stats.value?.popularItems ?? [])

const quickActions = computed(() => [
  { icon: Plus, label: t('menu.addItem'), to: { name: 'AdminDishCreate' }, color: '#e74c3c' },
  { icon: Food, label: t('routes.menuManagement'), to: { name: 'AdminMenuList' }, color: '#f39c12' },
  { icon: List, label: t('routes.orderManagement'), to: { name: 'AdminOrderList' }, color: '#3498db' },
  { icon: Grid, label: t('routes.tableManagement'), to: { name: 'AdminTableManage' }, color: '#27ae60' },
  { icon: Setting, label: t('routes.storeSettings'), to: { name: 'AdminStoreSettings' }, color: '#7c3aed' },
])

function openRoute(to: RouteLocationRaw) {
  router.push(to)
}
</script>

<template>
  <div class="dashboard admin-page">
    <AdminPageHeader :title="t('routes.dashboard')" :description="t('dashboard.overview')">
      <template #actions>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          :start-placeholder="t('dashboard.startDate')"
          :end-placeholder="t('dashboard.endDate')"
          :shortcuts="shortcuts"
          format="YYYY-MM-DD"
          :clearable="false"
          size="default"
          style="max-width: 300px"
        />
      </template>
    </AdminPageHeader>

    <el-row :gutter="12">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="never" class="stat stat-orders" v-loading="loading">
          <div class="stat-label">{{ t('dashboard.totalOrders') }}</div>
          <div class="stat-value">{{ totalOrders }}</div>
          <div class="stat-decor"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="never" class="stat stat-revenue" v-loading="loading">
          <div class="stat-label">{{ t('dashboard.paidRevenue') }}</div>
          <div class="stat-value">¥{{ paidRevenue.toFixed(2) }}</div>
          <div class="stat-decor"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="never" class="stat stat-pending" v-loading="loading">
          <div class="stat-label">{{ t('dashboard.pendingOrders') }}</div>
          <div class="stat-value">{{ pendingOrders }}</div>
          <div class="stat-decor"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="never" class="stat stat-tables" v-loading="loading">
          <div class="stat-label">{{ t('dashboard.tablesInUse') }}</div>
          <div class="stat-value">{{ tablesInUse }}</div>
          <div class="stat-decor"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="12" style="margin-top: 12px">
      <el-col :xs="24" :lg="14">
        <el-card shadow="never" class="card" v-loading="loading">
          <template #header>
            <div class="card-title">{{ t('dashboard.revenueChart') }}</div>
          </template>

          <div class="bars">
            <div v-for="d in revenueChart.days" :key="d.date" class="bar-item">
              <div class="bar" :style="{ height: `${(d.amount / revenueChart.max) * 120}px` }"></div>
              <div class="bar-label">{{ d.date }}</div>
              <div class="bar-val">¥{{ d.amount.toFixed(0) }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="10">
        <el-card shadow="never" class="card" v-loading="loading">
          <template #header>
            <div class="card-title">{{ t('dashboard.popularItems') }}</div>
          </template>

          <el-empty v-if="!popularItems.length" :description="t('common.noData')" />
          <div v-else class="rank">
            <div v-for="(d, idx) in popularItems" :key="d.dishId" class="rank-row">
              <div class="rank-idx" :class="{ gold: idx === 0, silver: idx === 1, bronze: idx === 2 }">{{ idx + 1 }}</div>
              <div class="rank-name">{{ localText(d.dishName) }}</div>
              <div class="rank-count">×{{ d.count }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="12" style="margin-top: 12px">
      <el-col :xs="24">
        <el-card shadow="never" class="card" v-loading="loading">
          <template #header>
            <div class="card-title">{{ t('common.action') }}</div>
          </template>

          <div class="quick-grid">
            <div
              v-for="a in quickActions"
              :key="a.label"
              class="quick-item hoverable-card"
              :style="{ '--accent': a.color }"
              role="button"
              tabindex="0"
              @click="openRoute(a.to)"
              @keydown.enter.prevent="openRoute(a.to)"
              @keydown.space.prevent="openRoute(a.to)"
            >
              <div class="quick-icon">
                <el-icon :size="18"><component :is="a.icon" /></el-icon>
              </div>
              <div class="quick-label">{{ a.label }}</div>
              <div class="quick-arrow">→</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped lang="scss">
.dashboard {
  width: 100%;
}

.stat {
  border-radius: 14px;
  position: relative;
  overflow: hidden;
  transition: transform var(--app-transition-base),
              box-shadow var(--app-transition-base);

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.10);
  }
}

.stat-decor {
  position: absolute;
  right: -16px;
  top: -16px;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  opacity: 0.15;
  pointer-events: none;
}

.stat-orders .stat-decor { background: #e74c3c; }
.stat-orders .stat-value { color: #e74c3c; }
.stat-revenue .stat-decor { background: #27ae60; }
.stat-revenue .stat-value { color: #27ae60; }
.stat-pending .stat-decor { background: #f39c12; }
.stat-pending .stat-value { color: #f39c12; }
.stat-tables .stat-decor { background: #3498db; }
.stat-tables .stat-value { color: #3498db; }

.stat-label {
  font-size: 12px;
  color: #909399;
  letter-spacing: 0.3px;
}

.stat-value {
  margin-top: 6px;
  font-size: 24px;
  font-weight: 900;
}

.card {
  border-radius: 14px;
  transition: box-shadow var(--app-transition-base);

  &:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  }
}

.card-title {
  font-weight: 700;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px;
}

.quick-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 12px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.85);
  border: 1px solid rgba(2, 6, 23, 0.06);
  cursor: pointer;
  user-select: none;
  min-width: 0;
}

.quick-icon {
  width: 34px;
  height: 34px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  background: linear-gradient(135deg, var(--accent), rgba(0, 0, 0, 0.25));
  flex: none;
}

.quick-label {
  font-weight: 800;
  color: #111827;
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.quick-arrow {
  margin-left: auto;
  color: rgba(17, 24, 39, 0.50);
  font-weight: 900;
}

.bars {
  height: 170px;
  display: flex;
  align-items: flex-end;
  gap: 10px;
  padding: 8px 4px;
}

.bar-item {
  flex: 1;
  text-align: center;
  min-width: 0;
}

.bar {
  width: 100%;
  background: linear-gradient(180deg, rgba(231, 76, 60, 0.35), rgba(231, 76, 60, 0.12));
  border: 1px solid rgba(231, 76, 60, 0.25);
  border-radius: 10px;
  min-height: 6px;
  transition: height 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.bar-label {
  margin-top: 6px;
  font-size: 12px;
  color: #909399;
}

.bar-val {
  margin-top: 2px;
  font-size: 12px;
  color: #606266;
}

.rank {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.rank-row {
  display: grid;
  grid-template-columns: 24px 1fr 60px;
  gap: 8px;
  align-items: center;
}

.rank-idx {
  font-weight: 900;
  color: #909399;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  font-size: 12px;

  &.gold {
    background: linear-gradient(135deg, #f9d423, #f39c12);
    color: #fff;
  }

  &.silver {
    background: linear-gradient(135deg, #bdc3c7, #95a5a6);
    color: #fff;
  }

  &.bronze {
    background: linear-gradient(135deg, #e67e22, #d35400);
    color: #fff;
  }
}

.rank-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #303133;
  font-weight: 600;
}

.rank-count {
  text-align: right;
  color: #909399;
  font-weight: 700;
}

@media (max-width: 1200px) {
  .quick-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 700px) {
  .quick-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
