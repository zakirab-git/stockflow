# StockFlow

Inventory and order API with a React admin UI.

## Quick start

### Option A — full stack (Docker)

```bash
docker compose up --build -d
```

Starts Postgres (applies `backend/migrations/001_init.sql` + `backend/scripts/seed.sql`), the Go API on `:8080`, and the UI on `:5173` (nginx proxies `/api` to the backend).

Schema changes after first boot: `docker compose down -v && docker compose up --build -d`.

### Option B — local API + UI (Postgres via Docker)

```bash
docker compose up -d postgres
```

#### Backend

```bash
cd backend
cp .env.example .env
go run .
```

| Variable | Example |
|----------|---------|
| `APP_NAME` | `stockflow` |
| `APP_ENV` | `development` |
| `APP_PORT` | `8080` |
| `DATABASE_HOST` | `localhost` |
| `DATABASE_PORT` | `5432` |
| `DATABASE_USER` | `stockflow` |
| `DATABASE_PASSWORD` | `stockflow` |
| `DATABASE_NAME` | `stockflow` |

#### Frontend

```bash
cd frontend/vite-project
npm install
npm run dev
```

Vite proxies `/api` → `http://127.0.0.1:8080` (see `vite.config.ts`).

## API examples

```bash
# List products
curl -s 'http://localhost:8080/api/products?page=1&limit=10' | jq

# Stock adjustment
curl -s -X POST 'http://localhost:8080/api/products/1/stock' \
  -H 'Content-Type: application/json' \
  -d '{"type":"out","qty":1,"note":"sample"}'

# Create order (409 if stock insufficient)
curl -s -X POST 'http://localhost:8080/api/orders' \
  -H 'Content-Type: application/json' \
  -d '{"customer":"Alice","items":[{"sku":"PE-010","qty":99}]}'

# Cancel order (restores stock)
curl -s -X POST 'http://localhost:8080/api/orders/1/cancel'
```

Error responses use `{ "error": { "code", "message"?, "available"?, "sku"? } }`.

## Tests

With Postgres up:

```bash
cd backend
STOCKFLOW_REQUIRE_DB=1 go test ./usecase/... -count=1
```

Optional `TEST_DATABASE_URL` overrides the default DSN. Without Postgres, tests skip unless `STOCKFLOW_REQUIRE_DB=1`.

## Design notes

- **last_movement_at**: derived in the products list query via a `LATERAL` subquery on `stock_movements` (no denormalized column).
- **Stock safety on orders**: `UPDATE products SET stock = stock - $qty WHERE id = $id AND stock >= $qty`; zero rows updated → `INSUFFICIENT_STOCK`. Order creation runs in a transaction with `SELECT … FOR UPDATE` on products, locking SKUs in sorted order.
- **Order status**: `CREATED` / `CANCELLED` (PDF contract). Product `PUT` updates name/price/active only; stock changes go through `POST …/stock`.
- **CORS vs proxy**: production builds may set `VITE_API_URL`; local dev uses the Vite proxy (`VITE_API_URL=/api`) to avoid CORS during review.

## Time spent

_Hours placeholder: ~[fill in]_

## What I'd improve next

- Idempotency keys on order create, list orders endpoint, server-side product filters, and a `go test -race` concurrency suite.