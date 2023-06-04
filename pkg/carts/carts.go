package carts

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Product struct {
	Name string `json:"name"`
	QTY  int64  `json:"qty"`
}
type ProductQty struct {
	Name string
	QTY  int64
}

func (s *Service) SaveToCart(ctx context.Context, id int64, products *[]Product) error {
	productQtyArr := make([]ProductQty, 0, len(*products))
	for _, product := range *products {
		productQtyArr = append(productQtyArr, ProductQty{Name: product.Name, QTY: product.QTY})
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO cart (customer_id, items)
		VALUES ($1, $2)
	`, id, productQtyArr)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
