package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction membuat transaksi baru beserta detail transaksi dalam satu database transaction.
//
// [ID]
// Fungsi ini akan:
// 1. Memvalidasi setiap produk (cek keberadaan & stok)
// 2. Mengunci row produk menggunakan FOR UPDATE (mencegah race condition)
// 3. Mengurangi stok produk
// 4. Menghitung total transaksi
// 5. Menyimpan data transaksi ke tabel `transactions`
// 6. Menyimpan detail transaksi ke tabel `transaction_details`
// 7. Melakukan commit jika semua proses berhasil, atau rollback jika terjadi error
//
// [EN]
// This function creates a new transaction along with its transaction details
// within a single database transaction.
//
// It will:
// 1. Validate each product (existence & stock availability)
// 2. Lock product rows using FOR UPDATE (to prevent race conditions)
// 3. Deduct product stock
// 4. Calculate total transaction amount
// 5. Insert transaction data into `transactions` table
// 6. Insert transaction details into `transaction_details` table
// 7. Commit if everything succeeds, or rollback if any error occurs
func (repo *TransactionRepository) CreateTransaction(
	items []models.CheckoutItem,
) (*models.Transaction, error) {

	// Memulai database transaction
	// Start database transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	// Flag untuk memastikan rollback hanya terjadi jika commit gagal
	// Flag to ensure rollback only happens if commit fails
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	// Total keseluruhan transaksi
	// Total transaction amount
	totalAmount := 0

	// Menyimpan detail transaksi sebelum di-insert ke database
	// Store transaction details before inserting into database
	details := make([]models.TransactionDetail, 0)

	// Loop setiap item checkout
	// Iterate through checkout items
	for _, item := range items {
		var (
			productID   int
			productName string
			price       int
			stock       int
		)

		// Mengambil data produk dan mengunci row (FOR UPDATE)
		// Fetch product data and lock row (FOR UPDATE)
		err := tx.QueryRow(
			"SELECT id, name, price, stock FROM products WHERE id = $1 FOR UPDATE",
			item.ProductID,
		).Scan(&productID, &productName, &price, &stock)

		// Jika produk tidak ditemukan
		// If product does not exist
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Validasi stok
		// Stock validation
		if stock < item.Quantity {
			return nil, fmt.Errorf(
				"stock not enough for product %s (available: %d, requested: %d)",
				productName, stock, item.Quantity,
			)
		}

		// Hitung subtotal per item
		// Calculate item subtotal
		subtotal := item.Quantity * price
		totalAmount += subtotal

		// Kurangi stok produk
		// Deduct product stock
		_, err = tx.Exec(
			"UPDATE products SET stock = stock - $1 WHERE id = $2",
			item.Quantity,
			productID,
		)
		if err != nil {
			return nil, err
		}

		// Tambahkan ke list transaction details
		// Append to transaction details list
		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Insert data transaksi utama
	// Insert main transaction data
	var transactionID int
	err = tx.QueryRow(
		"INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id",
		totalAmount,
	).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// Prepare statement untuk insert transaction details
	// Prepare statement for inserting transaction details
	stmt, err := tx.Prepare(`
		INSERT INTO transaction_details
		(transaction_id, product_id, quantity, subtotal)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Insert setiap detail transaksi
	// Insert each transaction detail
	for i := range details {
		details[i].TransactionID = transactionID

		_, err := stmt.Exec(
			transactionID,
			details[i].ProductID,
			details[i].Quantity,
			details[i].Subtotal,
		)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	committed = true

	// Return hasil transaksi
	// Return transaction result
	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
