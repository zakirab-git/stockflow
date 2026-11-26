import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"

import type { Movement } from "@/types/product"

interface MovementTableProps {
  movements: Movement[]
  loading?: boolean
}

export default function MovementTable({
  movements,
  loading,
}: MovementTableProps) {
  if (loading) {
    return (
      <p className="text-sm text-muted-foreground">Loading movements...</p>
    )
  }

  if (movements.length === 0) {
    return (
      <p className="text-sm text-muted-foreground">No movements recorded.</p>
    )
  }

  return (
    <div className="rounded-lg border bg-background">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Type</TableHead>
            <TableHead>Qty</TableHead>
            <TableHead>Note</TableHead>
            <TableHead>When</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {movements.map((m) => (
            <TableRow key={m.id}>
              <TableCell>
                <Badge variant={m.type === "in" ? "secondary" : "outline"}>
                  {m.type}
                </Badge>
              </TableCell>
              <TableCell>{m.qty}</TableCell>
              <TableCell className="max-w-[200px] truncate">
                {m.note || "—"}
              </TableCell>
              <TableCell>
                {new Date(m.created_at).toLocaleString("id-ID")}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}