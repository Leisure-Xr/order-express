<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useCartStore } from '@/stores/cart'
import EmptyState from '@/components/common/EmptyState.vue'
import CartItemRow from '@/components/customer/CartItemRow.vue'

const { t } = useI18n()
const router = useRouter()
const cartStore = useCartStore()

const formattedSubtotal = computed(() => `¥${cartStore.subtotal.toFixed(2)}`)

function goMenu() {
  router.push({ name: 'CustomerMenu' })
}

function goCheckout() {
  router.push({ name: 'CustomerCheckout' })
}

async function clearCart() {
  try {
    await ElMessageBox.confirm(t('cart.confirmClear'), t('common.confirm'), { type: 'warning' })
  } catch {
    return
  }
  cartStore.clearCart()
}
</script>

<template>
  <div class="cart-page">
    <div class="header">
      <div class="title">{{ t('cart.myCart') }}</div>
      <el-button v-if="!cartStore.isEmpty" text @click="clearCart">{{ t('cart.clearCart') }}</el-button>
    </div>

    <EmptyState
      v-if="cartStore.isEmpty"
      :title="t('cart.emptyCart')"
      :description="t('cart.goAddFood')"
    >
      <el-button type="primary" @click="goMenu">{{ t('routes.menu') }}</el-button>
    </EmptyState>

    <div v-else class="list">
      <CartItemRow
        v-for="item in cartStore.items"
        :key="item.key"
        :item="item"
        @update:quantity="({ key, quantity }) => cartStore.updateQuantity(key, quantity)"
        @remove="(key) => cartStore.removeItem(key)"
      />
    </div>

    <div v-if="!cartStore.isEmpty" class="footer">
      <div class="subtotal">
        {{ t('cart.subtotal') }}: <span class="price">{{ formattedSubtotal }}</span>
      </div>
      <el-button type="primary" size="large" class="checkout-btn" @click="goCheckout">
        {{ t('cart.goCheckout') }}
      </el-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.cart-page {
  padding: 12px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.title {
  font-size: 18px;
  font-weight: 800;
  color: #303133;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-bottom: 88px;
}

.footer {
  position: fixed;
  left: 0;
  right: 0;
  bottom: var(--customer-bottom-bar-height);
  margin: 0 auto;
  padding: 12px;
  max-width: 560px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 -1px 6px rgba(0, 0, 0, 0.06);
}

.subtotal {
  font-size: 14px;
  color: #606266;
}

.price {
  color: #e74c3c;
  font-weight: 900;
}

.checkout-btn {
  flex: none;
}
</style>
