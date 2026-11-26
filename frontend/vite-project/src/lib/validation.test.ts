import { describe, expect, it } from "vitest"
import { validateOrderLines, validateStockQty } from "./validation"

describe("validateOrderLines", () => {
  it("rejects empty lines", () => {
    expect(validateOrderLines([])).toBe("Add at least one line item")
  })

  it("rejects duplicate skus", () => {
    expect(
      validateOrderLines([
        { sku: "A", qty: 1 },
        { sku: "A", qty: 2 },
      ])
    ).toBe("Duplicate SKU in order: A")
  })
})

describe("validateStockQty", () => {
  it("requires positive integer", () => {
    expect(validateStockQty(0)).toBe("Quantity must be a positive integer")
    expect(validateStockQty(2.5)).toBe("Quantity must be a positive integer")
    expect(validateStockQty(3)).toBeNull()
  })
})