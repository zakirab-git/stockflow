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
import { Textarea } from "@/components/ui/textarea"
import { Alert, AlertDescription } from "@/components/ui/alert"

import { ApiError } from "@/services/api"
import { productService } from "@/services/productService"
import { validateStockQty } from "@/lib/validation"

import type { Product } from "@/types/product"

interface StockDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  product: Product | null
  onSaved: () => void
}

export default function StockDialog({
  open,
  onOpenChange,
  product,
  onSaved,
}: StockDialogProps) {
  const [type, setType] = useState<"in" | "out">("in")
  const [qty, setQty] = useState("1")
  const [note, setNote] = useState("")
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  useEffect(() => {
    if (open) {
      setType("in")
      setQty("1")
      setNote("")
      setError("")
    }
  }, [open, product])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!product) return

    const parsedQty = Number(qty)
    const validation = validateStockQty(parsedQty)
    if (validation) {
      setError(validation)
      return
    }

    try {
      setLoading(true)
      setError("")
      await productService.adjustStock(product.id, {
        type,
        qty: parsedQty,
        note: note.trim() || undefined,
      })
      onSaved()
      onOpenChange(false)
    } catch (err) {
      if (err instanceof ApiError && err.code === "INSUFFICIENT_STOCK") {
        setError(
          `Insufficient stock. Available: ${err.available ?? product.stock}`
        )
      } else {
        setError(err instanceof Error ? err.message : "Stock update failed")
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Adjust stock</DialogTitle>
          <DialogDescription>
            {product
              ? `${product.name} (${product.sku}) — current stock ${product.stock}`
              : ""}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label>Type</Label>
            <div className="flex gap-2">
              <Button
                type="button"
                variant={type === "in" ? "default" : "outline"}
                onClick={() => setType("in")}
              >
                Stock in
              </Button>
              <Button
                type="button"
                variant={type === "out" ? "default" : "outline"}
                onClick={() => setType("out")}
              >
                Stock out
              </Button>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="qty">Quantity</Label>
            <Input
              id="qty"
              type="number"
              min={1}
              step={1}
              value={qty}
              onChange={(e) => setQty(e.target.value)}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="note">Note</Label>
            <Textarea
              id="note"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              placeholder="Optional note"
              rows={2}
            />
          </div>

          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
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
            <Button type="submit" disabled={loading || !product}>
              {loading ? "Saving..." : "Save"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}