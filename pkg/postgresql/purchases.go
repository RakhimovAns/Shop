package postgresql

import (
	"context"
	"github.com/RakhimovAns/Shop/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type PurchaseService struct {
	pool *pgxpool.Pool
}

func NewPurchaseService(pool *pgxpool.Pool) *PurchaseService {
	return &PurchaseService{pool: pool}
}

func (s *PurchaseService) AddToPurchase(ctx context.Context, products []*types.Product, id int64) error {
	for _, product := range products {
		_, err := s.pool.Exec(ctx, `
			insert into purchases(customer_id, product_id, qty) VALUES ($1,$2,$3)
`, id, product.ID, product.QTY)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (s *PurchaseService) GetAllPurchase(ctx context.Context, id int64) ([]*types.Purchase, error) {
	items := make([]*types.Purchase, 0)
	rows, err := s.pool.Query(ctx, `
		select id,customer_id,product_id,qty,created from purchases where customer_id=$1
`, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		purchase := &types.Purchase{}
		err = rows.Scan(&purchase.ID, &purchase.CustomerId, &purchase.ProductId, &purchase.QTY, &purchase.Created)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, purchase)
	}
	return items, err
}
func (s *PurchaseService) DeletePurchase(ctx context.Context, id int64) error {
	_, err := s.pool.Exec(ctx, `
delete from purchases where customer_id=$1
`, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
