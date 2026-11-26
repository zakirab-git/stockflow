package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/zakirab-git/stockflow/apperrors"
	"github.com/zakirab-git/stockflow/dto"
	"github.com/zakirab-git/stockflow/repo"
	"github.com/zakirab-git/stockflow/usecase"
	"github.com/zakirab-git/stockflow/utils"
)

func testDB(t *testing.T) *sqlx.DB {
	t.Helper()

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		host := getenv("DATABASE_HOST", "localhost")
		port := getenv("DATABASE_PORT", "5432")
		user := getenv("DATABASE_USER", "heritsam")
		pass := getenv("DATABASE_PASSWORD", "local")
		name := getenv("DATABASE_NAME", "stockflow")
		url = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
	}

	db, err := sqlx.Connect("pgx", url)
	if err != nil {
		if os.Getenv("STOCKFLOW_REQUIRE_DB") == "1" {
			t.Fatalf("postgres required: %v", err)
		}
		t.Skipf("postgres not available: %v", err)
	}

	t.Cleanup(func() { db.Close() })
	return db
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func TestCreateOrder_InsufficientStockWritesNothing(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()

	txM := utils.NewTxM(db)
	orderRepo := repo.NewOrderRepo(db)
	uc := usecase.NewOrderUsecase(orderRepo, txM)

	okSKU := "TEST-OK-STOCK"
	shortSKU := "TEST-SHORT-STOCK"
	customer := "Tester-Insufficient"

	_, err := db.ExecContext(ctx, `
		INSERT INTO products (sku, name, price, stock, is_active)
		VALUES ($1, 'Ok', 1000, 50, true)
		ON CONFLICT (sku) DO UPDATE SET stock = 50
	`, okSKU)
	if err != nil {
		t.Fatalf("seed ok product: %v", err)
	}
	_, err = db.ExecContext(ctx, `
		INSERT INTO products (sku, name, price, stock, is_active)
		VALUES ($1, 'Short', 1000, 2, true)
		ON CONFLICT (sku) DO UPDATE SET stock = 2
	`, shortSKU)
	if err != nil {
		t.Fatalf("seed short product: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM sales_order_items WHERE sku IN ($1, $2)`, okSKU, shortSKU)
		_, _ = db.ExecContext(ctx, `DELETE FROM sales_orders WHERE customer = $1`, customer)
		_, _ = db.ExecContext(ctx, `DELETE FROM products WHERE sku IN ($1, $2)`, okSKU, shortSKU)
	})

	var beforeOk, beforeShort int
	if err := db.GetContext(ctx, &beforeOk, `SELECT stock FROM products WHERE sku = $1`, okSKU); err != nil {
		t.Fatalf("read ok stock: %v", err)
	}
	if err := db.GetContext(ctx, &beforeShort, `SELECT stock FROM products WHERE sku = $1`, shortSKU); err != nil {
		t.Fatalf("read short stock: %v", err)
	}

	_, _, _, err = uc.CreateOrder(ctx, dto.CreateOrderRequest{
		Customer: customer,
		Items: []dto.OrderItemRequest{
			{SKU: okSKU, Qty: 1},
			{SKU: shortSKU, Qty: 5},
		},
	})
	if !errors.Is(err, apperrors.ErrInsufficientStock) {
		t.Fatalf("expected INSUFFICIENT_STOCK, got %v", err)
	}

	var afterOk, afterShort int
	if err := db.GetContext(ctx, &afterOk, `SELECT stock FROM products WHERE sku = $1`, okSKU); err != nil {
		t.Fatalf("read ok stock after: %v", err)
	}
	if err := db.GetContext(ctx, &afterShort, `SELECT stock FROM products WHERE sku = $1`, shortSKU); err != nil {
		t.Fatalf("read short stock after: %v", err)
	}
	if afterOk != beforeOk || afterShort != beforeShort {
		t.Fatalf("stock changed: ok %d→%d short %d→%d", beforeOk, afterOk, beforeShort, afterShort)
	}

	var orderCount int
	if err := db.GetContext(ctx, &orderCount, `
		SELECT COUNT(*) FROM sales_orders WHERE customer = $1
	`, customer); err != nil {
		t.Fatalf("count orders: %v", err)
	}
	if orderCount != 0 {
		t.Fatalf("expected no orders, got %d", orderCount)
	}

	var itemCount int
	if err := db.GetContext(ctx, &itemCount, `
		SELECT COUNT(*) FROM sales_order_items WHERE sku IN ($1, $2)
	`, okSKU, shortSKU); err != nil {
		t.Fatalf("count items: %v", err)
	}
	if itemCount != 0 {
		t.Fatalf("expected no order items, got %d", itemCount)
	}
}

func TestUpdateStock_HappyPathCreatesMovement(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()

	txM := utils.NewTxM(db)
	productRepo := repo.NewProductRepo(db)
	uc := usecase.NewProductUsecase(productRepo, txM)

	sku := "TEST-STOCK-IN"
	var productID int
	err := db.GetContext(ctx, &productID, `
		INSERT INTO products (sku, name, price, stock, is_active)
		VALUES ($1, 'Stock In', 5000, 10, true)
		ON CONFLICT (sku) DO UPDATE SET stock = 10
		RETURNING id
	`, sku)
	if err != nil {
		t.Fatalf("seed product: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM stock_movements WHERE product_id = $1`, productID)
		_, _ = db.ExecContext(ctx, `DELETE FROM products WHERE sku = $1`, sku)
	})

	_, _, err = uc.UpdateStock(ctx, productID, dto.UpdateStockRequest{
		Type: "in",
		Qty:  3,
		Note: "restock",
	})
	if err != nil {
		t.Fatalf("update stock: %v", err)
	}

	var stock int
	if err := db.GetContext(ctx, &stock, `SELECT stock FROM products WHERE id = $1`, productID); err != nil {
		t.Fatalf("read stock: %v", err)
	}
	if stock != 13 {
		t.Fatalf("expected stock 13, got %d", stock)
	}

	var movementCount int
	if err := db.GetContext(ctx, &movementCount, `
		SELECT COUNT(*) FROM stock_movements
		WHERE product_id = $1 AND type = 'in' AND qty = 3
	`, productID); err != nil {
		t.Fatalf("count movements: %v", err)
	}
	if movementCount < 1 {
		t.Fatal("expected stock movement row")
	}
}
