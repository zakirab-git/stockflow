import { useCallback, useEffect, useMemo, useState } from "react"

import { Button } from "@/components/ui/button"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Input } from "@/components/ui/input"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"

import ProductTable from "@/components/products/ProductTable"
import ProductDialog from "@/components/products/ProductDialog"
import DeleteProductDialog from "@/components/products/DeleteProductDialog"
import StockDialog from "@/components/products/StockDialog"
import MovementTable from "@/components/products/MovementTable"

import { productService } from "@/services/productService"
import { formatRupiah } from "@/lib/format"
import { Plus, ChevronLeft, ChevronRight } from "lucide-react"

import type {
  Product,
  AddProductRequest,
  UpdateProductRequest,
} from "@/types/product"

type ProductFilter = "all" | "active" | "inactive" | "low"

const PAGE_SIZE = 10

export default function ProductsPage() {
  const [allProducts, setAllProducts] = useState<Product[]>([])
  const [searchInput, setSearchInput] = useState("")
  const [search, setSearch] = useState("")
  const [filter, setFilter] = useState<ProductFilter>("all")
  const [page, setPage] = useState(1)

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  const [dialogOpen, setDialogOpen] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [productToDelete, setProductToDelete] = useState<Product | null>(null)
  const [stockProduct, setStockProduct] = useState<Product | null>(null)
  const [stockOpen, setStockOpen] = useState(false)
  const [movementProduct, setMovementProduct] = useState<Product | null>(null)
  const [movementOpen, setMovementOpen] = useState(false)
  const [movements, setMovements] = useState<
    Awaited<ReturnType<typeof productService.getMovements>>
  >([])
  const [movementsLoading, setMovementsLoading] = useState(false)

  useEffect(() => {
    const timer = window.setTimeout(() => {
      setSearch(searchInput.trim())
      setPage(1)
    }, 350)
    return () => window.clearTimeout(timer)
  }, [searchInput])

  const fetchProducts = useCallback(async () => {
    try {
      setLoading(true)
      setError("")
      const { data } = await productService.getAllProducts(
        search || undefined
      )
      setAllProducts(data)
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to retrieve products."
      )
    } finally {
      setLoading(false)
    }
  }, [search])

  useEffect(() => {
    fetchProducts()
  }, [fetchProducts])

  const filtered = useMemo(() => {
    return allProducts.filter((p) => {
      if (filter === "active" && !p.is_active) return false
      if (filter === "inactive" && p.is_active) return false
      if (filter === "low" && p.stock > 10) return false
      return true
    })
  }, [allProducts, filter])

  const totalPages = Math.max(1, Math.ceil(filtered.length / PAGE_SIZE))
  const pageProducts = filtered.slice(
    (page - 1) * PAGE_SIZE,
    page * PAGE_SIZE
  )

  const stats = useMemo(() => {
    const active = allProducts.filter((p) => p.is_active).length
    const stockValue = allProducts.reduce(
      (sum, p) => sum + p.price * p.stock,
      0
    )
    const lowStock = allProducts.filter((p) => p.stock <= 10).length
    return {
      total: allProducts.length,
      active,
      stockValue,
      lowStock,
    }
  }, [allProducts])

  async function handleProductSubmit(
    formData: AddProductRequest | UpdateProductRequest
  ) {
    if (selectedProduct) {
      await productService.updateProduct(
        selectedProduct.id,
        formData as UpdateProductRequest
      )
    } else {
      await productService.addProduct(formData as AddProductRequest)
    }
    await fetchProducts()
  }

  async function confirmDelete() {
    if (!productToDelete) return
    await productService.deleteProduct(productToDelete.id)
    await fetchProducts()
    setProductToDelete(null)
  }

  async function openMovements(product: Product) {
    setMovementProduct(product)
    setMovementOpen(true)
    setMovementsLoading(true)
    try {
      const data = await productService.getMovements(product.id)
      setMovements(data)
    } finally {
      setMovementsLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold">Products</h2>
          <p className="text-sm text-muted-foreground">
            Manage inventory, stock movements, and pricing.
          </p>
        </div>
        <Button onClick={() => { setSelectedProduct(null); setDialogOpen(true) }}>
          <Plus className="mr-2 h-4 w-4" />
          Add Product
        </Button>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total products
            </CardTitle>
          </CardHeader>
          <CardContent className="text-2xl font-bold">{stats.total}</CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Active
            </CardTitle>
          </CardHeader>
          <CardContent className="text-2xl font-bold">{stats.active}</CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Stock value
            </CardTitle>
          </CardHeader>
          <CardContent className="text-2xl font-bold">
            {formatRupiah(stats.stockValue)}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Low stock (≤10)
            </CardTitle>
          </CardHeader>
          <CardContent className="text-2xl font-bold">{stats.lowStock}</CardContent>
        </Card>
      </div>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <Input
          className="max-w-sm"
          placeholder="Search name or SKU..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
        />
        <div className="flex flex-wrap gap-2">
          {(
            [
              ["all", "All"],
              ["active", "Active"],
              ["inactive", "Inactive"],
              ["low", "Low stock"],
            ] as const
          ).map(([key, label]) => (
            <Button
              key={key}
              size="sm"
              variant={filter === key ? "default" : "outline"}
              onClick={() => {
                setFilter(key)
                setPage(1)
              }}
            >
              {label}
            </Button>
          ))}
        </div>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription className="flex items-center justify-between gap-4">
            <span>{error}</span>
            <Button size="sm" variant="outline" onClick={fetchProducts}>
              Retry
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {loading ? (
        <p className="text-sm text-muted-foreground">Loading products...</p>
      ) : filtered.length === 0 ? (
        <p className="text-sm text-muted-foreground">No products match filters.</p>
      ) : (
        <ProductTable
          products={pageProducts}
          onEdit={(p) => { setSelectedProduct(p); setDialogOpen(true) }}
          onDelete={(id) => {
            const product = allProducts.find((p) => p.id === id)
            if (product) {
              setProductToDelete(product)
              setDeleteDialogOpen(true)
            }
          }}
          onStock={(p) => { setStockProduct(p); setStockOpen(true) }}
          onMovements={openMovements}
        />
      )}

      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Showing {pageProducts.length} of {filtered.length} filtered (
          {allProducts.length} loaded)
        </p>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            disabled={page <= 1 || loading}
            onClick={() => setPage((p) => p - 1)}
          >
            <ChevronLeft className="h-4 w-4" />
            Previous
          </Button>
          <span className="text-sm">
            Page {page} of {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={page >= totalPages || loading}
            onClick={() => setPage((p) => p + 1)}
          >
            Next
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <ProductDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        product={selectedProduct}
        onSubmit={handleProductSubmit}
      />

      <DeleteProductDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        productName={productToDelete?.name ?? ""}
        onConfirm={confirmDelete}
      />

      <StockDialog
        open={stockOpen}
        onOpenChange={setStockOpen}
        product={stockProduct}
        onSaved={fetchProducts}
      />

      <Dialog open={movementOpen} onOpenChange={setMovementOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>
              Movement history — {movementProduct?.name}
            </DialogTitle>
          </DialogHeader>
          <MovementTable movements={movements} loading={movementsLoading} />
        </DialogContent>
      </Dialog>
    </div>
  )
}