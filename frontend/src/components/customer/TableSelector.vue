<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import type { Table } from '@/types'
import { getPublicTablesApi } from '@/api/table'

const props = withDefaults(
  defineProps<{
    modelValue: string
    disabled?: boolean
  }>(),
  {
    disabled: false,
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const { t } = useI18n()

const loading = ref(false)
const tables = ref<Table[]>([])

const availableTables = computed(() => {
  const list = tables.value
    .filter((tb) => tb.status === 'available')
    .sort((a, b) => a.number.localeCompare(b.number))

  // If a table number comes from QR query but isn't available (occupied/reserved),
  // still show it for selection.
  if (props.modelValue && !tables.value.find((tb) => tb.number === props.modelValue)) {
    list.unshift({
      id: `custom-${props.modelValue}`,
      number: props.modelValue,
      seats: 0,
      status: 'available',
    })
  }

  return list
})

onMounted(async () => {
  loading.value = true
  try {
    const res = await getPublicTablesApi()
    if (res.code !== 200) {
      ElMessage.error(res.message || t('common.fail'))
      return
    }
    tables.value = res.data ?? []
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <el-select
    :model-value="modelValue"
    :disabled="disabled"
    :loading="loading"
    filterable
    style="width: 220px"
    :placeholder="t('order.selectTable')"
    @update:model-value="(v: string) => emit('update:modelValue', v)"
  >
    <el-option
      v-for="tb in availableTables"
      :key="tb.id"
      :value="tb.number"
      :label="tb.area ? `${tb.number} - ${tb.area}` : tb.number"
    />
  </el-select>
</template>
