package service

import (
	"context"
	"github.com/RakhimovAns/Shop/pkg/postgresql"
	"github.com/RakhimovAns/Shop/types"
)

type CustomerService struct {
	service postgresql.CustomerService
}

func NewCustomerService(service *postgresql.CustomerService) *CustomerService {
	return &CustomerService{service: *service}
}

func (s *CustomerService) Register(ctx context.Context, customer *types.Customer) error {
	return s.service.Register(ctx, customer)
}

func (s *CustomerService) Login(ctx context.Context, login string, password string) (string, error) {
	return s.service.Login(ctx, login, password)
}

func (s *CustomerService) Delete(ctx context.Context, id int64) error {
	return s.service.Delete(ctx, id)
}

func (s *CustomerService) GetByID(ctx context.Context, id int64) (error, *types.Customer) {
	return s.service.GetByID(ctx, id)
}
func (s *CustomerService) ChangeBalance(ctx context.Context, id int64, sum int64) error {
	return s.service.ChangeBalance(ctx, id, sum)
}
func (s *CustomerService) DepositBalance(ctx context.Context, id int64, sum int64) error {
	return s.service.DepositBalance(ctx, id, sum)
}
