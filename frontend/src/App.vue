<script setup lang="ts">
import { computed } from 'vue'
import { RouterView } from 'vue-router'
import { useRoute } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'
import CustomerLayout from '@/layouts/CustomerLayout.vue'
import BlankLayout from '@/layouts/BlankLayout.vue'

type LayoutName = 'admin' | 'customer' | 'blank'

const route = useRoute()

const layoutName = computed<LayoutName>(() => (route.meta.layout as LayoutName) || 'blank')

const layoutComponent = computed(() => {
  switch (layoutName.value) {
    case 'admin':
      return AdminLayout
    case 'customer':
      return CustomerLayout
    default:
      return BlankLayout
  }
})

const blankShowBack = computed(() => {
  if (layoutName.value !== 'blank') return undefined
  if (typeof route.meta.showBack === 'boolean') return route.meta.showBack
  if (route.name === 'AdminLogin') return false
  return true
})

const layoutProps = computed(() => (layoutName.value === 'blank' ? { showBack: blankShowBack.value } : {}))
</script>

<template>
  <component :is="layoutComponent" v-bind="layoutProps">
    <RouterView v-slot="{ Component }">
      <Transition name="page-fade" mode="out-in">
        <component :is="Component" />
      </Transition>
    </RouterView>
  </component>
</template>

<style>
#app {
  width: 100%;
  min-height: 100vh;
}
</style>
