package service

import (
	"context"
	"github.com/RakhimovAns/Shop/pkg/postgresql"
	"github.com/RakhimovAns/Shop/types"
)

type CartService struct {
	service *postgresql.CartService
}

func NewCartService(service *postgresql.CartService) *CartService {
	return &CartService{service: service}
}
func (s *CartService) GetCartID(ctx context.Context, id int64) (error, int64) {
	return s.service.GetCartID(ctx, id)
}
func (s *CartService) SaveToCart(ctx context.Context, id int64, products *[]types.Product) error {
	return s.service.SaveToCart(ctx, id, products)
}

func (s *CartService) DeleteProducts(ctx context.Context, id int64, products *[]types.Product) error {
	return s.service.DeleteProducts(ctx, id, products)
}

func (s *CartService) ChangeQTY(ctx context.Context, id int64, products *[]types.Product) error {
	return s.service.ChangeQTY(ctx, id, products)
}

func (s *CartService) GetCartBYID(ctx context.Context, id int64) ([]*types.Product, error) {
	return s.service.GetCartBYID(ctx, id)
}

func (s *CartService) GetSum(ctx context.Context, products []*types.Product) (int64, error) {
	return s.service.GetSum(ctx, products)
}
func (s *CartService) DeleteCart(ctx context.Context, ID int64) error {
	return s.service.DeleteCart(ctx, ID)
}
