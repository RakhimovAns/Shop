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
	ID  int64 `json:"id"`
	QTY int64 `json:"qty"`
}

func (s *Service) CreateCart(ctx context.Context, id int64) (error, int64) {
	_, err := s.pool.Exec(ctx, `
insert into carts(customer_id) values ($1) on conflict do nothing 
`, id)
	if err != nil {
		log.Println(err)
		return err, -1
	}
	var ID int64
	err = s.pool.QueryRow(ctx, `
select id from carts where customer_id=$1
`, id).Scan(&ID)
	if err != nil {
		log.Println(err)
		return err, -1
	}

	return nil, ID
}
func (s *Service) SaveToCart(ctx context.Context, id int64, products *[]Product) error {
	for _, product := range *products {
		_, err := s.pool.Exec(ctx, `
insert into carts_items(cart_id, product_id, count) VALUES ($1,$2,$3) on conflict(cart_id,product_id) do update set count=carts_items.count+excluded.count
`, id, product.ID, product.QTY)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
