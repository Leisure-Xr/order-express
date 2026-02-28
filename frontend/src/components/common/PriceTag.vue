<script setup lang="ts">
import { computed } from 'vue'
import { formatPrice } from '@/utils/format'

const props = withDefaults(
  defineProps<{
    price: number
    originalPrice?: number
    size?: 'sm' | 'md' | 'lg'
  }>(),
  {
    size: 'md',
  },
)

const formattedPrice = computed(() => formatPrice(props.price))
const formattedOriginalPrice = computed(() =>
  props.originalPrice ? formatPrice(props.originalPrice) : '',
)
const hasDiscount = computed(
  () => props.originalPrice !== undefined && props.originalPrice > props.price,
)
</script>

<template>
  <span class="price-tag" :class="[`price-tag--${size}`]">
    <span class="price-current">{{ formattedPrice }}</span>
    <span v-if="hasDiscount" class="price-original">{{ formattedOriginalPrice }}</span>
  </span>
</template>

<style scoped lang="scss">
.price-tag {
  display: inline-flex;
  align-items: baseline;
  gap: 4px;
}

.price-current {
  color: #e74c3c;
  font-weight: 700;
}

.price-original {
  color: #999999;
  text-decoration: line-through;
}

// Sizes
.price-tag--sm {
  .price-current {
    font-size: 14px;
  }

  .price-original {
    font-size: 12px;
  }
}

.price-tag--md {
  .price-current {
    font-size: 18px;
  }

  .price-original {
    font-size: 14px;
  }
}

.price-tag--lg {
  .price-current {
    font-size: 24px;
  }

  .price-original {
    font-size: 16px;
  }
}
</style>
