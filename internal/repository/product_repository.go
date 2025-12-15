package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dominikuswilly/nofu-be_product/internal/entity"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id int64) (*entity.Product, error)
	GetAll(ctx context.Context) ([]*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id int64) error
}

// postgresProductRepository implements ProductRepository for PostgreSQL
type postgresProductRepository struct {
	db *sql.DB
}

// NewPostgresProductRepository creates a new postgresProductRepository
func NewPostgresProductRepository(db *sql.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Create(ctx context.Context, product *entity.Product) error {
	query := `
		INSERT INTO products (name, description, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *postgresProductRepository) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	query := `
		SELECT id, name, description, price, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	product := &entity.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Return nil if found nothing, let usecase handle 404
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	return product, nil
}

func (r *postgresProductRepository) GetAll(ctx context.Context) ([]*entity.Product, error) {
	query := `
		SELECT id, name, description, price, created_at, updated_at
		FROM products
		ORDER BY id ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		product := &entity.Product{}
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *postgresProductRepository) Update(ctx context.Context, product *entity.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, updated_at = $4
		WHERE id = $5
	`
	product.UpdatedAt = time.Now()
	res, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.UpdatedAt,
		product.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *postgresProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}
