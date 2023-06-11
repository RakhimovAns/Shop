package service

import (
	"context"
	"github.com/RakhimovAns/Shop/pkg/postgresql"
	"github.com/RakhimovAns/Shop/types"
)

type PurchaseService struct {
	service *postgresql.PurchaseService
}

func NewPurchaseService(service *postgresql.PurchaseService) *PurchaseService {
	return &PurchaseService{service: service}
}

func (s *PurchaseService) AddToPurchase(ctx context.Context, products []*types.Product, id int64) error {
	return s.service.AddToPurchase(ctx, products, id)
}

func (s *PurchaseService) GetAllPurchase(ctx context.Context, id int64) ([]*types.Purchase, error) {
	return s.service.GetAllPurchase(ctx, id)
}

func (s *PurchaseService) DeletePurchase(ctx context.Context, id int64) error {
	return s.service.DeletePurchase(ctx, id)
}
