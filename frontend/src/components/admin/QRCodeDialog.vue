<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { generateQRCodeDataUrl } from '@/utils/qrcode'

const props = defineProps<{
  modelValue: boolean
  tableNumber: string
  link: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const { t } = useI18n()

const visible = computed({
  get: () => props.modelValue,
  set: (v: boolean) => emit('update:modelValue', v),
})

const loading = ref(false)
const qrDataUrl = ref<string>('')

watch(
  () => props.modelValue,
  async (v) => {
    if (!v) return
    loading.value = true
    try {
      qrDataUrl.value = await generateQRCodeDataUrl(props.link)
    } catch (e: any) {
      ElMessage.error(e?.message || t('common.fail'))
      qrDataUrl.value = ''
    } finally {
      loading.value = false
    }
  },
)

function download() {
  if (!qrDataUrl.value) return
  const a = document.createElement('a')
  a.href = qrDataUrl.value
  a.download = `table-${props.tableNumber}.png`
  a.click()
}
</script>

<template>
  <el-dialog v-model="visible" :title="`${t('table.scanToOrder')} - ${tableNumber}`" width="420px">
    <div class="qr-body">
      <div v-loading="loading" class="qr-box">
        <img v-if="qrDataUrl" :src="qrDataUrl" alt="QR" />
      </div>
      <div class="qr-link">
        <el-input :model-value="link" readonly />
      </div>
      <div class="qr-actions">
        <el-button @click="download">{{ t('table.downloadQRCode') }}</el-button>
        <el-button type="primary" @click="visible = false">{{ t('common.confirm') }}</el-button>
      </div>
    </div>
  </el-dialog>
</template>

<style scoped lang="scss">
.qr-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.qr-box {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  border-radius: 12px;
  padding: 12px;
  min-height: 320px;
}

.qr-box img {
  width: 320px;
  height: 320px;
  border-radius: 8px;
}

.qr-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>

