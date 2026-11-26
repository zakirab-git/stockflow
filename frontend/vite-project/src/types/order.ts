export interface OrderItem {
  sku: string
  qty: number
  unit_price: number
}

export interface Order {
  id: number
  customer: string
  items: OrderItem[]
  total: number
  status: string
  created_at: string
}