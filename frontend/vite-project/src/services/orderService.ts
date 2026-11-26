import { apiRequest } from "./api"
import type { Order } from "@/types/order"

export const orderService = {
  createOrder(data: {
    customer: string
    items: {
      sku: string
      qty: number
    }[]
  }) {
    return apiRequest<Order>("/orders", {
      method: "POST",
      body: JSON.stringify(data),
    })
  },

  getOrder(id: number) {
    return apiRequest<Order>(`/orders/${id}`)
  },

  cancelOrder(id: number) {
    return apiRequest<Order>(`/orders/${id}/cancel`, {
      method: "POST",
    })
  },
}