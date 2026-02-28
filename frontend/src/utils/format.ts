export function formatPrice(price: number): string {
  return `¥${price.toFixed(2)}`
}

export function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatOrderId(id: string): string {
  return `#${id.slice(-6).toUpperCase()}`
}

export function generateCartItemKey(
  dishId: string,
  selectedOptions: { optionName: string; valueName: string }[],
): string {
  const optionStr = selectedOptions
    .map((o) => `${o.optionName}:${o.valueName}`)
    .sort()
    .join('|')
  return `${dishId}_${optionStr || 'default'}`
}
