import { ref, onUnmounted } from 'vue'

export function useOrderPolling() {
  const timer = ref<ReturnType<typeof setInterval> | null>(null)
  const isPolling = ref(false)

  function startPolling(callback: () => void | Promise<void>, interval = 5000) {
    stopPolling()
    isPolling.value = true
    // Execute immediately on start
    callback()
    timer.value = setInterval(() => {
      callback()
    }, interval)
  }

  function stopPolling() {
    if (timer.value !== null) {
      clearInterval(timer.value)
      timer.value = null
    }
    isPolling.value = false
  }

  onUnmounted(() => {
    stopPolling()
  })

  return {
    isPolling,
    startPolling,
    stopPolling,
  }
}
