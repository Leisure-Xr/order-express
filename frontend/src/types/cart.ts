export interface CartItem {
  key: string
  dishId: string
  dishName: { zh: string; en: string }
  image: string
  unitPrice: number
  quantity: number
  selectedOptions: { optionName: string; valueName: string; priceAdjustment: number }[]
}
