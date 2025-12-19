package server

import (
	"context"
	"net/http"

	"github.com/dominikuswilly/nofu-be_product/internal/config"
	"github.com/dominikuswilly/nofu-be_product/internal/handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func NewServer(cfg *config.Config, handler *handler.ProductHandler, logger *zap.Logger) *Server {
	router := gin.Default()

	// Global middleware
	router.Use(gin.Recovery())
	// Logger middleware already included in Default, but we can customize if needed

	// Register routes
	api := router.Group("/api/product")
	handler.RegisterRoutes(api)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	return &Server{
		httpServer: srv,
		logger:     logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting server", zap.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
