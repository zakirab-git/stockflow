export interface Product {
  id: number
  sku: string
  name: string
  price: number
  stock: number
  is_active: boolean
  last_movement_at?: string | null
}

export interface ProductListResponse {
  data: Product[]
  total: number
  page: number
  limit: number
  total_pages?: number
}

export interface AddProductRequest {
  sku: string
  name: string
  price: number
  stock: number
  is_active: boolean
}

export interface UpdateProductRequest {
  name: string
  price: number
  is_active: boolean
}

export interface Movement {
  id: number
  product_id: number
  type: "in" | "out"
  qty: number
  note?: string
  created_at: string
}

export interface UpdateStockRequest {
  type: "in" | "out"
  qty: number
  note?: string
}