import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

import { orderService } from "@/services/orderService"
import { ApiError } from "@/services/api"
import { validateOrderLines } from "@/lib/validation"
import { formatRupiah } from "@/lib/format"
import { Plus, Trash2 } from "lucide-react"

import type { Order } from "@/types/order"

interface OrdersPageProps {
  onViewOrder: (id: number) => void
  refreshToken?: number
}

type LineItem = { sku: string; qty: string }

export default function OrdersPage({ onViewOrder }: OrdersPageProps) {
  const [customer, setCustomer] = useState("")
  const [items, setItems] = useState<LineItem[]>([{ sku: "", qty: "1" }])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const [conflictSku, setConflictSku] = useState<string | null>(null)
  const [created, setCreated] = useState<Order | null>(null)

  function updateItem(index: number, patch: Partial<LineItem>) {
    setItems((prev) =>
      prev.map((row, i) => (i === index ? { ...row, ...patch } : row))
    )
  }

  function addRow() {
    setItems((prev) => [...prev, { sku: "", qty: "1" }])
  }

  function removeRow(index: number) {
    setItems((prev) => prev.filter((_, i) => i !== index))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError("")
    setConflictSku(null)
    setCreated(null)

    const parsed = items.map((row) => ({
      sku: row.sku.trim(),
      qty: Number(row.qty),
    }))

    const validation = validateOrderLines(parsed)
    if (validation) {
      setError(validation)
      return
    }

    if (!customer.trim()) {
      setError("Customer name is required")
      return
    }

    try {
      setLoading(true)
      const order = await orderService.createOrder({
        customer: customer.trim(),
        items: parsed,
      })
      setCreated(order)
      setCustomer("")
      setItems([{ sku: "", qty: "1" }])
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.code === "INSUFFICIENT_STOCK") {
          setConflictSku(err.sku ?? null)
          setError(
            `Insufficient stock for ${err.sku ?? "item"} (available: ${err.available ?? "?"})`
          )
        } else {
          setError(err.message)
        }
      } else {
        setError(err instanceof Error ? err.message : "Failed to create order")
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold">Orders</h2>
        <p className="text-sm text-muted-foreground">
          Create sales orders; stock is decremented in the same transaction.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Create order</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2 max-w-md">
              <Label htmlFor="customer">Customer</Label>
              <Input
                id="customer"
                value={customer}
                onChange={(e) => setCustomer(e.target.value)}
                placeholder="Customer name"
              />
            </div>

            <div className="space-y-3">
              <Label>Line items</Label>
              {items.map((row, index) => {
                const isConflict =
                  conflictSku !== null &&
                  row.sku.trim().toLowerCase() === conflictSku.toLowerCase()
                return (
                <div
                  key={index}
                  className={`flex flex-wrap items-end gap-2 rounded-md p-2 ${
                    isConflict ? "bg-destructive/10 ring-2 ring-destructive" : ""
                  }`}
                >
                  <div className="space-y-1">
                    <Label className="text-xs text-muted-foreground">SKU</Label>
                    <Input
                      className="w-36"
                      value={row.sku}
                      onChange={(e) =>
                        updateItem(index, { sku: e.target.value })
                      }
                      placeholder="SKU"
                      aria-invalid={isConflict}
                    />
                  </div>
                  <div className="space-y-1">
                    <Label className="text-xs text-muted-foreground">Qty</Label>
                    <Input
                      className="w-24"
                      type="number"
                      min={1}
                      step={1}
                      value={row.qty}
                      onChange={(e) =>
                        updateItem(index, { qty: e.target.value })
                      }
                      aria-invalid={isConflict}
                    />
                  </div>
                  <Button
                    type="button"
                    variant="outline"
                    size="icon"
                    disabled={items.length === 1}
                    onClick={() => removeRow(index)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              )})}
              <Button type="button" variant="outline" size="sm" onClick={addRow}>
                <Plus className="mr-2 h-4 w-4" />
                Add line
              </Button>
            </div>

            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <Button type="submit" disabled={loading}>
              {loading ? "Creating..." : "Create order"}
            </Button>
          </form>
        </CardContent>
      </Card>

      {created && (
        <Alert>
          <AlertDescription className="space-y-2">
            <p>
              Order #{created.id} created for {created.customer} — total{" "}
              {formatRupiah(created.total)}.
            </p>
            <Button size="sm" onClick={() => onViewOrder(created.id)}>
              View order
            </Button>
          </AlertDescription>
        </Alert>
      )}
    </div>
  )
}