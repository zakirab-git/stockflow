export type OrderLineInput = { sku: string; qty: number }

export function validateOrderLines(
  lines: OrderLineInput[]
): string | null {
  if (lines.length === 0) {
    return "Add at least one line item"
  }

  const seen = new Set<string>()

  for (const line of lines) {
    const sku = line.sku.trim()
    if (!sku) {
      return "Each line needs a SKU"
    }
    if (seen.has(sku)) {
      return `Duplicate SKU in order: ${sku}`
    }
    seen.add(sku)

    if (!Number.isInteger(line.qty) || line.qty < 1) {
      return `Quantity for ${sku} must be a positive integer`
    }
  }

  return null
}

export function validateStockQty(qty: number): string | null {
  if (!Number.isInteger(qty) || qty < 1) {
    return "Quantity must be a positive integer"
  }
  return null
}