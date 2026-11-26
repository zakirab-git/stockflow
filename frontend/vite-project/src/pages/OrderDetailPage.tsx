import { useCallback, useEffect, useState } from "react"

import { Button } from "@/components/ui/button"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Badge } from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"

import { orderService } from "@/services/orderService"
import { ApiError } from "@/services/api"
import { formatRupiah } from "@/lib/format"
import { ArrowLeft } from "lucide-react"

import type { Order } from "@/types/order"

interface OrderDetailPageProps {
  orderId: number
  onBack: () => void
  onCancelled?: () => void
}

export default function OrderDetailPage({
  orderId,
  onBack,
  onCancelled,
}: OrderDetailPageProps) {
  const [order, setOrder] = useState<Order | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [cancelOpen, setCancelOpen] = useState(false)
  const [cancelling, setCancelling] = useState(false)

  const load = useCallback(async () => {
    try {
      setLoading(true)
      setError("")
      const data = await orderService.getOrder(orderId)
      setOrder(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load order")
    } finally {
      setLoading(false)
    }
  }, [orderId])

  useEffect(() => {
    load()
  }, [load])

  async function confirmCancel() {
    try {
      setCancelling(true)
      const updated = await orderService.cancelOrder(orderId)
      setOrder(updated)
      onCancelled?.()
    } catch (err) {
      setError(
        err instanceof ApiError
          ? err.message
          : err instanceof Error
            ? err.message
            : "Cancel failed"
      )
    } finally {
      setCancelling(false)
      setCancelOpen(false)
    }
  }

  if (loading) {
    return <p className="text-sm text-muted-foreground">Loading order...</p>
  }

  if (!order) {
    return (
      <div className="space-y-4">
        <Button variant="outline" onClick={onBack}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back
        </Button>
        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between gap-4">
        <div>
          <Button variant="outline" size="sm" onClick={onBack} className="mb-3">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to orders
          </Button>
          <h2 className="text-2xl font-bold">Order #{order.id}</h2>
          <p className="text-sm text-muted-foreground">
            {order.customer} · {new Date(order.created_at).toLocaleString("id-ID")}
          </p>
        </div>
        <Badge variant={order.status === "CREATED" ? "secondary" : "outline"}>
          {order.status}
        </Badge>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className="rounded-lg border bg-background p-4">
        <p className="text-lg font-semibold">
          Total {formatRupiah(order.total)}
        </p>
      </div>

      <div className="rounded-lg border bg-background">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>SKU</TableHead>
              <TableHead>Qty</TableHead>
              <TableHead>Unit price</TableHead>
              <TableHead>Subtotal</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {order.items.map((item) => (
              <TableRow key={`${item.sku}-${item.qty}`}>
                <TableCell>{item.sku}</TableCell>
                <TableCell>{item.qty}</TableCell>
                <TableCell>{formatRupiah(item.unit_price)}</TableCell>
                <TableCell>
                  {formatRupiah(item.unit_price * item.qty)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {order.status === "CREATED" && (
        <Button variant="destructive" onClick={() => setCancelOpen(true)}>
          Cancel order
        </Button>
      )}

      <AlertDialog open={cancelOpen} onOpenChange={setCancelOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Cancel this order?</AlertDialogTitle>
            <AlertDialogDescription>
              Stock for all line items will be restored.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={cancelling}>Keep order</AlertDialogCancel>
            <AlertDialogAction disabled={cancelling} onClick={confirmCancel}>
              {cancelling ? "Cancelling..." : "Cancel order"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}