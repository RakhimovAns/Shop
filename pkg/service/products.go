package service

import (
	"context"
	"github.com/RakhimovAns/Shop/pkg/postgresql"
	"github.com/RakhimovAns/Shop/types"
)

type ProductService struct {
	service *postgresql.ProductService
}

func NewProductService(service *postgresql.ProductService) *ProductService {
	return &ProductService{service: service}
}

func (s *ProductService) AllActiveProducts(ctx context.Context) ([]*types.Product, error) {
	return s.service.AllActiveProducts(ctx)
}

func (s *ProductService) AllCategories(ctx context.Context) ([]*string, error) {
	return s.service.AllCategories(ctx)
}
func (s *ProductService) GetByCategory(ctx context.Context, category string) ([]*types.Product, error) {
	return s.service.GetByCategory(ctx, category)
}

func (s *ProductService) Search(ctx context.Context, Key string) ([]*types.Product, error) {
	return s.service.Search(ctx, Key)
}
