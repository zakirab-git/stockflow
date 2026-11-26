import { apiRequest } from "./api"

import type {
  Product,
  ProductListResponse,
  AddProductRequest,
  UpdateProductRequest,
  Movement,
  UpdateStockRequest,
} from "@/types/product"

export interface ProductQuery {
  page?: number
  limit?: number
  search?: string
}

function buildParams(query: ProductQuery) {
  const params = new URLSearchParams()
  params.set("page", String(query.page ?? 1))
  params.set("limit", String(query.limit ?? 10))
  if (query.search) {
    params.set("search", query.search)
  }
  return params
}

export const productService = {
  getProducts(query: ProductQuery = {}) {
    return apiRequest<ProductListResponse>(
      `/products?${buildParams(query).toString()}`
    )
  },

  async getAllProducts(search?: string) {
    const limit = 20
    let page = 1
    const items: Product[] = []
    let total = 0

    for (;;) {
      const result = await this.getProducts({ page, limit, search })
      total = result.total
      const pageItems = result.data ?? []
      items.push(...pageItems)
      if (items.length >= total || pageItems.length === 0) {
        break
      }
      page += 1
    }

    return { data: items, total }
  },

  addProduct(data: AddProductRequest) {
    return apiRequest<Product>("/products", {
      method: "POST",
      body: JSON.stringify(data),
    })
  },

  updateProduct(id: number, data: UpdateProductRequest) {
    return apiRequest<Product>(`/products/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    })
  },

  deleteProduct(id: number) {
    return apiRequest<void>(`/products/${id}`, {
      method: "DELETE",
    })
  },

  adjustStock(id: number, data: UpdateStockRequest) {
    return apiRequest<Movement>(`/products/${id}/stock`, {
      method: "POST",
      body: JSON.stringify(data),
    })
  },

  async getMovements(id: number) {
    const res = await apiRequest<{ data: Movement[] }>(
      `/products/${id}/movements`
    )
    return res.data ?? []
  },
}