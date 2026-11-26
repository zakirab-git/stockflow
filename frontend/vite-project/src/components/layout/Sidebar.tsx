import {
  Package,
  ShoppingCart,
  Boxes,
} from "lucide-react"

interface SidebarProps {
  activePage: string
  onNavigate: (page: string) => void
}

export default function Sidebar({
  activePage,
  onNavigate,
}: SidebarProps) {
  const navigation = [
    {
      name: "Products",
      icon: Package,
    },
    {
      name: "Orders",
      icon: ShoppingCart,
    },
  ]

  return (
    <aside className="w-64 border-r bg-background min-h-screen p-4">
      <div className="flex items-center gap-2 mb-10 px-2">
        <Boxes className="h-6 w-6" />
        <h1 className="text-xl font-bold">
          Inventory App
        </h1>
      </div>

      <nav className="space-y-2">
        {navigation.map((item) => {
          const Icon = item.icon
          const active = activePage === item.name

          return (
            <button
              key={item.name}
              onClick={() => onNavigate(item.name)}
              className={`flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors ${
                active
                  ? "bg-primary text-primary-foreground"
                  : "hover:bg-muted"
              }`}
            >
              <Icon className="h-4 w-4" />
              {item.name}
            </button>
          )
        })}
      </nav>
    </aside>
  )
}