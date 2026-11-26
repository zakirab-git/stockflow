import { useState } from "react"

import Sidebar from "@/components/layout/Sidebar"
import ProductsPage from "@/pages/ProductsPage"
import OrdersPage from "@/pages/OrdersPage"
import OrderDetailPage from "@/pages/OrderDetailPage"

type NavPage = "Products" | "Orders" | "OrderDetail"

function App() {
  const [activePage, setActivePage] = useState<NavPage>("Products")
  const [orderId, setOrderId] = useState<number | null>(null)
  const [productsRefresh, setProductsRefresh] = useState(0)

  function handleNavigate(page: string) {
    if (page === "Products" || page === "Orders") {
      setActivePage(page)
      setOrderId(null)
    }
  }

  function viewOrder(id: number) {
    setOrderId(id)
    setActivePage("OrderDetail")
  }

  return (
    <div className="flex min-h-screen">
      <Sidebar
        activePage={activePage === "OrderDetail" ? "Orders" : activePage}
        onNavigate={handleNavigate}
      />

      <main className="flex-1 bg-muted/30 p-8">
        <div className="mx-auto max-w-6xl">
          {activePage === "Products" && (
            <ProductsPage key={productsRefresh} />
          )}

          {activePage === "Orders" && (
            <OrdersPage onViewOrder={viewOrder} refreshToken={productsRefresh} />
          )}

          {activePage === "OrderDetail" && orderId !== null && (
            <OrderDetailPage
              orderId={orderId}
              onBack={() => setActivePage("Orders")}
              onCancelled={() => {
                setProductsRefresh((n) => n + 1)
                setActivePage("Products")
              }}
            />
          )}
        </div>
      </main>
    </div>
  )
}

export default App