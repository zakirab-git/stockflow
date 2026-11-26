# StockFlow

Small inventory/order API (Go + Postgres) with a React admin UI.

## Running it

Easiest way, everything in Docker:

```bash
docker compose up --build -d

API is on `:8080`, UI on `:5173`. Schema and seed data load on first boot. If you change the schema, reset with `docker compose down -v` and start again.
```

The `.env.example` has the DB settings (host, port, user, password, name all `localhost` / `5432` / `stockflow`). Vite proxies `/api` to `:8080`.

## Tests

```bash
cd backend
STOCKFLOW_REQUIRE_DB=1 go test ./usecase/... -count=1
```

i also included postman collection file to test the endpoints

Needs Postgres running. Without `STOCKFLOW_REQUIRE_DB=1` the tests just skip. You can point at another DB with `TEST_DATABASE_URL`.

Time spent
~12 hours

## Next up

i wished i spent more times on frontend. the things i missed are: Idempotency keys for order creation, a list-orders endpoint, server-side product filters, and a concurrency test with `-race`.