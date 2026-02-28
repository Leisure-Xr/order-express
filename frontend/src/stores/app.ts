import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const locale = ref<'zh-CN' | 'en'>(
    (localStorage.getItem('locale') as 'zh-CN' | 'en') || 'zh-CN',
  )
  const isLoading = ref(false)
  const sidebarCollapsed = ref(false)
  const currentTableId = ref<string | null>(
    localStorage.getItem('tableId') || null,
  )

  function setLocale(newLocale: 'zh-CN' | 'en') {
    locale.value = newLocale
    localStorage.setItem('locale', newLocale)
  }

  function setTableFromQuery(tableId: string | null) {
    currentTableId.value = tableId
    if (tableId) {
      localStorage.setItem('tableId', tableId)
    } else {
      localStorage.removeItem('tableId')
    }
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  return {
    locale,
    isLoading,
    sidebarCollapsed,
    currentTableId,
    setLocale,
    setTableFromQuery,
    toggleSidebar,
  }
})
