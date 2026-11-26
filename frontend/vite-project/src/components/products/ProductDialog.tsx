
import { useEffect, useState } from "react"

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"

import type {
  Product,
  AddProductRequest,
  UpdateProductRequest,
} from "@/types/product"

interface ProductDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  product: Product | null
  onSubmit: (
    data: AddProductRequest | UpdateProductRequest
  ) => Promise<void>
}

export default function ProductDialog({
  open,
  onOpenChange,
  product,
  onSubmit,
}: ProductDialogProps) {
  const [sku, setSku] = useState("")
  const [name, setName] = useState("")
  const [price, setPrice] = useState("0")
  const [stock, setStock] = useState("0")
  const [isActive, setIsActive] = useState(true)

  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  const isEditing = product !== null

  useEffect(() => {
    if (open) {
      setSku(product?.sku ?? "")
      setName(product?.name ?? "")
      setPrice(String(product?.price ?? 0))
      setStock(String(product?.stock ?? 0))
      setIsActive(product?.is_active ?? true)
      setError("")
    }
  }, [open, product])

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()

    try {
      setLoading(true)
      setError("")

      if (isEditing) {
        const data: UpdateProductRequest = {
          name,
          price: Number(price),
          is_active: isActive,
        }

        await onSubmit(data)
      } else {
        const data: AddProductRequest = {
          sku,
          name,
          price: Number(price),
          stock: Number(stock),
          is_active: isActive,
        }

        await onSubmit(data)
      }

      onOpenChange(false)
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to save product."
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            {isEditing ? "Edit Product" : "Add Product"}
          </DialogTitle>

          <DialogDescription>
            {isEditing
              ? "Update the product information."
              : "Enter the details of the new product."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          {!isEditing && (
            <div className="space-y-2">
              <Label htmlFor="sku">SKU</Label>
              <Input
                id="sku"
                value={sku}
                onChange={(e) => setSku(e.target.value)}
                placeholder="e.g. PROD-001"
                required
              />
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="name">Product Name</Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Product name"
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="price">Price</Label>
            <Input
              id="price"
              type="number"
              min="0"
              step="1"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              required
            />
          </div>

          {!isEditing && (
            <div className="space-y-2">
              <Label htmlFor="stock">Stock</Label>
              <Input
                id="stock"
                type="number"
                min="0"
                step="1"
                value={stock}
                onChange={(e) => setStock(e.target.value)}
                required
              />
            </div>
          )}

          <div className="flex items-center gap-2">
            <Checkbox
              id="is_active"
              checked={isActive}
              onCheckedChange={(checked) =>
                setIsActive(checked === true)
              }
            />
            <Label htmlFor="is_active">Active product</Label>
          </div>

          {error && (
            <p className="text-sm text-destructive">{error}</p>
          )}

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={loading}
            >
              Cancel
            </Button>

            <Button type="submit" disabled={loading}>
              {loading
                ? "Saving..."
                : isEditing
                  ? "Save Changes"
                  : "Add Product"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
