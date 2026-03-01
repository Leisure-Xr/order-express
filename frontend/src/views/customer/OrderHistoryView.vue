<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useCartStore } from '@/stores/cart'
import { useMenuStore } from '@/stores/menu'
import { useLocaleText } from '@/composables/useLocaleText'
import { formatDate } from '@/utils/format'
import { getOrderHistoryIds } from '@/utils/orderHistory'
import { getOrderByIdApi } from '@/api/order'
import type { Order, OrderStatus } from '@/types'
import EmptyState from '@/components/common/EmptyState.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()
const menuStore = useMenuStore()
const { localText } = useLocaleText()

const loading = ref(false)

const orders = ref<Order[]>([])

type OrderQuickFilter = 'all' | 'unpaid' | 'to_ship' | 'to_receive' | 'after_sales'

const activeFilter = computed<OrderQuickFilter>(() => {
  const value = route.query.filter
  if (typeof value !== 'string') return 'all'
  if (value === 'unpaid' || value === 'to_ship' || value === 'to_receive' || value === 'after_sales') {
    return value
  }
  return 'all'
})

const filterOptions = computed(() => [
  { label: t('common.all'), value: 'all' },
  { label: t('profile.toPay'), value: 'unpaid' },
  { label: t('profile.toShip'), value: 'to_ship' },
  { label: t('profile.toReceive'), value: 'to_receive' },
  { label: t('profile.afterSales'), value: 'after_sales' },
])

function setFilter(next: string | number) {
  const value = next as OrderQuickFilter
  if (value === activeFilter.value) return

  const query = { ...route.query }
  delete query.filter

  router.replace({
    name: 'CustomerOrderHistory',
    query: value === 'all' ? query : { ...query, filter: value },
  })
}

function matchesFilter(order: Order, filter: Exclude<OrderQuickFilter, 'all'>): boolean {
  if (filter === 'unpaid') return order.payment.status === 'unpaid'
  if (filter === 'to_ship') return order.status === 'confirmed' || order.status === 'preparing' || order.status === 'ready'
  if (filter === 'to_receive') return order.status === 'delivered'
  return order.payment.status === 'refunded' || order.status === 'cancelled'
}

const displayOrders = computed(() => {
  const filter = activeFilter.value
  if (filter === 'all') return orders.value
  return orders.value.filter((o) => matchesFilter(o, filter))
})

function statusTagType(status: OrderStatus): 'info' | 'success' | 'warning' | 'danger' {
  if (status === 'pending') return 'warning'
  if (status === 'confirmed' || status === 'preparing' || status === 'ready') return 'info'
  if (status === 'completed' || status === 'delivered') return 'success'
  return 'danger'
}

onMounted(async () => {
  loading.value = true
  try {
    const ids = getOrderHistoryIds()
    if (!ids.length) {
      orders.value = []
      return
    }

    const resolved = await Promise.all(
      ids.map(async (id) => {
        const res = await getOrderByIdApi(id)
        return res.code === 200 ? res.data : null
      }),
    )
    orders.value = resolved.filter((o): o is Order => !!o)
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  } finally {
    loading.value = false
  }
})

function track(order: Order) {
  router.push({ name: 'CustomerOrderTracking', params: { id: order.id } })
}

async function reorder(order: Order) {
  try {
    if (!menuStore.dishes.length) {
      await menuStore.fetchDishes()
    }

    for (const item of order.items) {
      let dish = menuStore.dishes.find((d) => d.id === item.dishId) ?? null
      if (!dish) {
        dish = await menuStore.fetchDishById(item.dishId)
      }
      if (!dish) continue

      cartStore.addItem(dish, item.quantity, item.selectedOptions)
    }

    ElMessage.success(t('common.success'))
    router.push({ name: 'CustomerCart' })
  } catch (e: any) {
    ElMessage.error(e?.message || t('common.fail'))
  }
}
</script>

<template>
  <div class="orders-page">
    <div class="title">{{ t('routes.orderHistory') }}</div>

    <div class="filters">
      <el-segmented
        :model-value="activeFilter"
        :options="filterOptions"
        @change="setFilter"
      />
    </div>

    <el-skeleton v-if="loading && !orders.length" :rows="8" animated />

    <EmptyState v-else-if="!displayOrders.length" :title="t('common.noData')" />

    <div v-else class="list">
      <el-card v-for="o in displayOrders" :key="o.id" shadow="never" class="order-card">
        <div class="top">
          <div class="order-no">{{ o.orderNumber }}</div>
          <el-tag :type="statusTagType(o.status)">{{ t(`order.status.${o.status}`) }}</el-tag>
        </div>

        <div class="meta">
          <span>{{ formatDate(o.createdAt) }}</span>
          <span class="dot">·</span>
          <span>
            {{
              o.type === 'dine_in'
                ? t('order.dineIn')
                : o.type === 'takeaway'
                  ? t('order.takeout')
                  : t('order.pickup')
            }}
          </span>
          <span v-if="o.tableId">
            <span class="dot">·</span>
            <span>{{ t('order.tableNumber') }}: {{ o.tableId }}</span>
          </span>
        </div>

        <div class="items">
          <div v-for="(it, idx) in o.items.slice(0, 3)" :key="idx" class="item-line">
            {{ localText(it.dishName) }} × {{ it.quantity }}
          </div>
          <div v-if="o.items.length > 3" class="more">+{{ o.items.length - 3 }}</div>
        </div>

        <div class="bottom">
          <div class="amount">¥{{ o.total.toFixed(2) }}</div>
          <div class="btns">
            <el-button size="small" @click="track(o)">{{ t('payment.trackOrder') }}</el-button>
            <el-button type="primary" size="small" @click="reorder(o)">{{ t('common.add') }}</el-button>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<style scoped lang="scss">
.orders-page {
  padding: 12px;
}

.title {
  font-size: 18px;
  font-weight: 900;
  color: #303133;
  margin-bottom: 12px;
}

.filters {
  margin-bottom: 12px;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-bottom: 12px;

  > :nth-child(1) { animation: fadeSlideUp 0.3s ease 0.05s both; }
  > :nth-child(2) { animation: fadeSlideUp 0.3s ease 0.10s both; }
  > :nth-child(3) { animation: fadeSlideUp 0.3s ease 0.15s both; }
  > :nth-child(4) { animation: fadeSlideUp 0.3s ease 0.20s both; }
  > :nth-child(5) { animation: fadeSlideUp 0.3s ease 0.25s both; }
}

.order-card {
  border-radius: 14px;
  border: 1px solid rgba(0, 0, 0, 0.03);
  transition: transform var(--app-transition-base),
              box-shadow var(--app-transition-base);

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.10);
  }
}

.top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.order-no {
  font-weight: 900;
  color: #303133;
}

.meta {
  font-size: 12px;
  color: #909399;
}

.dot {
  margin: 0 6px;
}

.items {
  margin-top: 10px;
  color: #606266;
  font-size: 13px;
}

.item-line {
  line-height: 1.6;
}

.more {
  color: #909399;
  font-size: 12px;
  margin-top: 2px;
}

.bottom {
  margin-top: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.amount {
  background: linear-gradient(135deg, #e74c3c, #f0574a);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 900;
  font-size: 16px;
}

.btns {
  display: flex;
  gap: 8px;
}
</style>
