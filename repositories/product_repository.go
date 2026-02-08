package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
	"strings"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := `SELECT
        p.id,
        p.name,
        p.price,
        p.stock,
        c.id AS category_id,
        c.name AS category_name
    FROM products p
    JOIN categories c ON p.category_id = c.id`
	var args []interface{}
	if strings.TrimSpace(name) != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryId, &p.CategoryName)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	err := repo.checkCategoryExists(product.CategoryId)
	if err != nil {
		return err
	}
	query := "INSERT INTO products (name, price, stock,category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err = repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryId).Scan(&product.ID)
	return err
}

func (repo *ProductRepository) checkCategoryExists(id int) error {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1) AS category_exists`
	result := false
	err := repo.db.QueryRow(query, id).Scan(&result)
	if err == sql.ErrNoRows || !result {
		return errors.New("category tidak ditemukan")
	}
	if err != nil {
		return err
	}
	return nil
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `SELECT
        p.id,
        p.name,
        p.price,
        p.stock,
        c.id AS category_id,
        c.name AS category_name
    FROM products p
    JOIN categories c ON p.category_id = c.id WHERE p.id = $1`

	var p models.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryId, &p.CategoryName)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	err := repo.checkCategoryExists(product.CategoryId)
	if err != nil {
		return err
	}

	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryId, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}
