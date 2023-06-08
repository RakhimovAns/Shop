package purchase

import (
	"context"
	"github.com/RakhimovAns/Shop/pkg/carts"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Purchase struct {
	ID         int64     `json:"id"`
	CustomerId int64     `json:"customer_id"`
	ProductId  int64     `json:"product_id"`
	QTY        int64     `json:"qty"`
	Created    time.Time `json:"created"`
}

func (s *Service) AddToPurchase(ctx context.Context, products []*carts.Product, id int64) error {
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

func (s *Service) GetAllPurchase(ctx context.Context, id int64) ([]*Purchase, error) {
	items := make([]*Purchase, 0)
	rows, err := s.pool.Query(ctx, `
select id,customer_id,product_id,qty,created from purchases where customer_id=$1
`, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		purchase := &Purchase{}
		err = rows.Scan(&purchase.ID, &purchase.CustomerId, &purchase.ProductId, &purchase.QTY, &purchase.Created)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, purchase)
	}
	return items, err
}
