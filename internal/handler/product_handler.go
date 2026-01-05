package handler

import (
	"net/http"

	"github.com/dominikuswilly/nofu-be_product/internal/dto"
	"github.com/dominikuswilly/nofu-be_product/internal/middleware"
	"github.com/dominikuswilly/nofu-be_product/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductHandler struct {
	usecase        usecase.ProductUsecase
	logger         *zap.Logger
	authServiceURL string
}

func NewProductHandler(usecase usecase.ProductUsecase, logger *zap.Logger, authServiceURL string) *ProductHandler {
	return &ProductHandler{
		usecase:        usecase,
		logger:         logger,
		authServiceURL: authServiceURL,
	}
}

func (h *ProductHandler) RegisterRoutes(r *gin.RouterGroup) {
	products := r.Group("/products")
	products.Use(middleware.AuthMiddleware(h.authServiceURL))
	{
		products.POST("", h.CreateProduct)
		products.GET("", h.GetAllProducts)
		products.GET("/:id", h.GetProductByID)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usecase.CreateProduct(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"responseCode":    "201",
		"responseMessage": "success",
		"data":            res,
	})
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	res, err := h.usecase.GetAllProducts(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to fetch products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"responseCode":    "200",
		"responseMessage": "success",
		"data":            res,
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.usecase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to fetch product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}
	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usecase.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		h.logger.Error("Failed to update product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := h.usecase.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		// Basic check for not found if repository returned specific error,
		// but here we just assume 500 or 404 based on string match or robust error handling
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		h.logger.Error("Failed to delete product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.Status(http.StatusNoContent)
}
