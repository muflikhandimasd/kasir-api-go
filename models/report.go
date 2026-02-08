package models

type BestSellerProduct struct {
	Name     string `json:"nama"`
	Quantity int64  `json:"qty_terjual"`
}

type SalesSummary struct {
	TotalRevenue   int                `json:"total_revenue"`
	TotalTransaksi int                `json:"total_transaksi"`
	ProdukTerlaris *BestSellerProduct `json:"produk_terlaris"`
}
