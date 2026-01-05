package usecase

import (
	"context"
	"time"

	"github.com/dominikuswilly/nofu-be_product/internal/dto"
	"github.com/dominikuswilly/nofu-be_product/internal/entity"
	"github.com/dominikuswilly/nofu-be_product/internal/repository"
	"github.com/google/uuid"
)

// ProductUsecase defines the business logic interface
type ProductUsecase interface {
	CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProductByID(ctx context.Context, id string) (*dto.ProductResponse, error)
	GetAllProducts(ctx context.Context) ([]*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id string) error
}

type productUsecase struct {
	repo repository.ProductRepository
}

// NewProductUsecase creates a new productUsecase
func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (u *productUsecase) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	createdBy := "orang"

	newID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	active := int16(1) // default value is 1
	if req.Active != nil {
		if *req.Active {
			active = 1
		} else {
			active = 0
		}
	}

	product := &entity.Product{
		ID:          newID.String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    req.Currency,
		Url:         req.Url,
		Stock:       *req.Stock,
		Active:      active,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}

	if err := u.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return toProductResponse(product), nil
}

func (u *productUsecase) GetProductByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil // Or return a specific ErrNotFound
	}
	return toProductResponse(product), nil
}

func (u *productUsecase) GetAllProducts(ctx context.Context) ([]*dto.ProductResponse, error) {
	products, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = toProductResponse(p)
	}
	return responses, nil
}

func (u *productUsecase) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// First check if exists
	existingProduct, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingProduct == nil {
		return nil, nil
	}

	// Update fields if present
	if req.Name != nil {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil {
		existingProduct.Description = *req.Description
	}
	if req.Price != nil {
		existingProduct.Price = *req.Price
	}
	if req.Stock != nil {
		existingProduct.Stock = *req.Stock
	}
	if req.Active != nil {
		if *req.Active {
			existingProduct.Active = 1
		} else {
			existingProduct.Active = 0
		}
	}

	if err := u.repo.Update(ctx, existingProduct); err != nil {
		return nil, err
	}

	return toProductResponse(existingProduct), nil
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func toProductResponse(p *entity.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Currency:    p.Currency,
		Url:         p.Url,
		CreatedBy:   p.CreatedBy,
		Stock:       p.Stock,
		Active:      p.Active == 1,
	}
}
