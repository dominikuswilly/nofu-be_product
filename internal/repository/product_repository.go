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
	GetByID(ctx context.Context, id string) (*entity.Product, error)
	GetAll(ctx context.Context) ([]*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id string) error
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
		INSERT INTO product_master (c_id, c_nm, c_description, d_price, c_currency, c_url, c_created_by, i_stock, i_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	err := r.db.QueryRowContext(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Currency,
		product.Url,
		product.CreatedBy,
		product.Stock,
		product.Active,
	).Err()

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *postgresProductRepository) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	query := `
		SELECT c_id, c_nm, c_description, d_price, c_currency, c_url, c_created_by, ts_created_at, ts_updated_at, i_stock, i_active
		FROM product_master
		WHERE c_id = $1
	`
	product := &entity.Product{}
	var createdAt, updatedAt sql.NullTime
	// var createdBy sql.NullString // Removed, scanning directly into product.CreatedBy
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Currency,
		&product.Url,
		&product.CreatedBy, // Scan directly into *string
		&createdAt,
		&updatedAt,
		&product.Stock,
		&product.Active,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Return nil if found nothing, let usecase handle 404
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}

	// product.CreatedBy = createdBy.String // Removed
	if createdAt.Valid {
		product.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		product.UpdatedAt = updatedAt.Time
	}
	return product, nil
}

func (r *postgresProductRepository) GetAll(ctx context.Context) ([]*entity.Product, error) {
	query := `
		SELECT c_id, c_nm, c_description, d_price, c_currency, c_url, c_created_by, ts_created_at, ts_updated_at, i_stock, i_active
		FROM product_master
		ORDER BY c_id ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		product := &entity.Product{}
		var createdAt, updatedAt sql.NullTime
		// var createdBy sql.NullString // Removed
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Currency,
			&product.Url,
			&product.CreatedBy, // Scan directly into *string
			&createdAt,
			&updatedAt,
			&product.Stock,
			&product.Active,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		// product.CreatedBy = createdBy.String // Removed
		if createdAt.Valid {
			product.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			product.UpdatedAt = updatedAt.Time
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *postgresProductRepository) Update(ctx context.Context, product *entity.Product) error {
	query := `
		UPDATE product_master
		SET c_nm = $1, c_description = $2, d_price = $3, i_stock = $4, ts_updated_at = $5, i_active = $6
		WHERE c_id = $7
	`
	product.UpdatedAt = time.Now()
	res, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.UpdatedAt,
		product.Active,
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

func (r *postgresProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM product_master WHERE c_id = $1`
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
