<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useOrderStore } from '@/stores/order'
import { useTableStore } from '@/stores/table'
import { useLocaleText } from '@/composables/useLocaleText'
import type { Order } from '@/types'
import { formatDate } from '@/utils/format'

const { t } = useI18n()
const { localText } = useLocaleText()
const orderStore = useOrderStore()
const tableStore = useTableStore()

const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([
      orderStore.fetchOrders({ status: 'all', page: 1, pageSize: 200 }),
      tableStore.fetchTables(),
    ])
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    loading.value = false
  }
})

const baseDate = computed(() => {
  const times = orderStore.orders.map((o) => new Date(o.createdAt).getTime())
  const max = times.length ? Math.max(...times) : Date.now()
  return new Date(max)
})

function sameDay(a: Date, b: Date): boolean {
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate()
}

const todayOrders = computed(() => {
  const today = baseDate.value
  return orderStore.orders.filter((o) => sameDay(new Date(o.createdAt), today)).length
})

const todayRevenue = computed(() => {
  const today = baseDate.value
  return orderStore.orders
    .filter((o) => sameDay(new Date(o.createdAt), today) && o.payment.status === 'paid')
    .reduce((sum, o) => sum + o.total, 0)
})

const pendingOrders = computed(() => orderStore.orders.filter((o) => o.status === 'pending').length)
const tablesInUse = computed(() => tableStore.occupiedTables.length)

const recentOrders = computed(() => orderStore.orders.slice(0, 6))

const popularDishes = computed(() => {
  const map = new Map<string, { name: { zh: string; en: string }; count: number }>()
  for (const o of orderStore.orders) {
    for (const it of o.items) {
      const cur = map.get(it.dishId)
      if (cur) cur.count += it.quantity
      else map.set(it.dishId, { name: it.dishName, count: it.quantity })
    }
  }
  return [...map.entries()]
    .map(([dishId, v]) => ({ dishId, ...v }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 6)
})

const revenueLastWeek = computed(() => {
  const end = new Date(baseDate.value)
  const days: { date: string; amount: number }[] = []
  for (let i = 6; i >= 0; i--) {
    const d = new Date(end)
    d.setDate(d.getDate() - i)
    const key = d.toISOString().slice(0, 10)
    const amount = orderStore.orders
      .filter((o) => o.payment.status === 'paid' && o.createdAt.slice(0, 10) === key)
      .reduce((sum, o) => sum + o.total, 0)
    days.push({ date: key.slice(5), amount })
  }

  const max = Math.max(...days.map((d) => d.amount), 1)
  return { days, max }
})
</script>

<template>
  <div class="dashboard">
    <div class="title">{{ t('routes.dashboard') }}</div>

    <el-row :gutter="12">
      <el-col :span="6">
        <el-card shadow="never" class="stat" v-loading="loading">
          <div class="stat-label">{{ t('dashboard.todayOrders') }}</div>
          <div class="stat-value">{{ todayOrders }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat" v-loading="loading">
          <div class="stat-label">{{ t('dashboard.todayRevenue') }}</div>
          <div class="stat-value">¥{{ todayRevenue.toFixed(2) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat" v-loading="loading">
          <div class="stat-label">{{ t('dashboard.pendingOrders') }}</div>
          <div class="stat-value">{{ pendingOrders }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat" v-loading="loading">
          <div class="stat-label">{{ t('dashboard.tablesInUse') }}</div>
          <div class="stat-value">{{ tablesInUse }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="12" style="margin-top: 12px">
      <el-col :span="14">
        <el-card shadow="never" class="card" v-loading="loading">
          <template #header>
            <div class="card-title">{{ t('dashboard.revenueLastWeek') }}</div>
          </template>

          <div class="bars">
            <div v-for="d in revenueLastWeek.days" :key="d.date" class="bar-item">
              <div class="bar" :style="{ height: `${(d.amount / revenueLastWeek.max) * 120}px` }"></div>
              <div class="bar-label">{{ d.date }}</div>
              <div class="bar-val">¥{{ d.amount.toFixed(0) }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="10">
        <el-card shadow="never" class="card" v-loading="loading">
          <template #header>
            <div class="card-title">{{ t('dashboard.popularItems') }}</div>
          </template>

          <div class="rank">
            <div v-for="(d, idx) in popularDishes" :key="d.dishId" class="rank-row">
              <div class="rank-idx">{{ idx + 1 }}</div>
              <div class="rank-name">{{ localText(d.name) }}</div>
              <div class="rank-count">×{{ d.count }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="12" style="margin-top: 12px">
      <el-col :span="24">
        <el-card shadow="never" class="card" v-loading="loading">
          <template #header>
            <div class="card-title">{{ t('dashboard.recentOrders') }}</div>
          </template>

          <el-table :data="recentOrders" stripe>
            <el-table-column label="Order" min-width="200">
              <template #default="{ row }">
                <div class="order-no">{{ row.orderNumber }}</div>
                <div class="order-time">{{ formatDate(row.createdAt) }}</div>
              </template>
            </el-table-column>
            <el-table-column :label="t('order.orderType')" width="120">
              <template #default="{ row }">
                {{
                  row.type === 'dine_in'
                    ? t('order.dineIn')
                    : row.type === 'takeaway'
                      ? t('order.takeout')
                      : t('order.pickup')
                }}
              </template>
            </el-table-column>
            <el-table-column :label="t('common.status')" width="140">
              <template #default="{ row }">
                {{ t(`order.status.${row.status}`) }}
              </template>
            </el-table-column>
            <el-table-column :label="t('payment.amount')" width="140">
              <template #default="{ row }">¥{{ row.total.toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped lang="scss">
.dashboard {
  max-width: 1400px;
}

.title {
  font-size: 18px;
  font-weight: 800;
  color: #303133;
  margin-bottom: 12px;
}

.stat {
  border-radius: 12px;
}

.stat-label {
  font-size: 12px;
  color: #909399;
}

.stat-value {
  margin-top: 6px;
  font-size: 22px;
  font-weight: 900;
  color: #303133;
}

.card {
  border-radius: 12px;
}

.card-title {
  font-weight: 700;
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
}

.bar {
  width: 100%;
  background: rgba(231, 76, 60, 0.22);
  border: 1px solid rgba(231, 76, 60, 0.35);
  border-radius: 10px;
  min-height: 6px;
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
  color: #e74c3c;
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

.order-no {
  font-weight: 800;
  color: #303133;
}

.order-time {
  font-size: 12px;
  color: #909399;
}

@media (max-width: 1200px) {
  :deep(.el-col-6) {
    width: 50%;
  }
}

@media (max-width: 800px) {
  :deep(.el-col-6) {
    width: 100%;
  }
}
</style>
