package repositories

import (
	"database/sql"
	"kasir-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetSalesSummary mendapatkan ringkasan penjualan berdasarkan rentang tanggal.
//
// [ID]
// Jika startDate & endDate bernilai nil, maka otomatis mengambil data hari ini.
// Jika keduanya diisi, maka akan mengambil data sesuai rentang tanggal.
//
// [EN]
// If startDate & endDate are nil, it returns today's sales summary.
// If both are provided, it returns sales summary for the given date range.
func (repo *ReportRepository) GetSalesSummary(
	startDate *time.Time,
	endDate *time.Time,
) (*models.SalesSummary, error) {

	var summary models.SalesSummary

	// =========================
	// 1. Total Revenue & Transaksi
	// =========================
	baseQuery := `
		SELECT
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COUNT(id) AS total_transaksi
		FROM transactions
		WHERE created_at BETWEEN $1 AND $2
	`

	var start, end time.Time

	if startDate == nil || endDate == nil {
		// Hari ini (00:00:00 - 23:59:59)
		now := time.Now()
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.Add(24 * time.Hour)
	} else {
		start = *startDate
		end = *endDate
	}

	err := repo.db.QueryRow(baseQuery, start, end).
		Scan(&summary.TotalRevenue, &summary.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// =========================
	// 2. Produk Terlaris
	// =========================
	var (
		name sql.NullString
		qty  sql.NullInt64
	)
	bestSellerQuery := `
		SELECT
			p.name,
			SUM(td.quantity) AS qty_terjual
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`

	err = repo.db.QueryRow(bestSellerQuery, start, end).
		Scan(
			&name,
			&qty,
		)

	// Jika tidak ada transaksi sama sekali
	if err == sql.ErrNoRows {
		return &summary, nil
	}

	// Jika query mengembalikan row tapi kolom NULL
	if !name.Valid || !qty.Valid {
		return &summary, nil
	}

	if err != nil {
		return nil, err
	}

	summary.ProdukTerlaris = &models.BestSellerProduct{
		Name:     name.String,
		Quantity: qty.Int64,
	}

	return &summary, nil
}
