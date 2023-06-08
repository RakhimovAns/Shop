package carts

import (
	"context"
	"errors"
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

var ErrNosuch = errors.New("more than you have")

func (s *Service) GetCartID(ctx context.Context, id int64) (error, int64) {
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

func (s *Service) DeleteProducts(ctx context.Context, id int64, products *[]Product) error {
	for _, product := range *products {
		_, err := s.pool.Exec(ctx, `
delete from carts_items where cart_id=$1 and product_id=$2
`, id, product.ID)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (s *Service) ChangeQTY(ctx context.Context, id int64, products *[]Product) error {
	for _, product := range *products {
		var count int64
		err := s.pool.QueryRow(ctx, `
select count from carts_items where cart_id=$1 and product_id=$2
`, id, product.ID).Scan(&count)
		if err != nil {
			log.Println(err)
			return err
		}
		if product.QTY == count {
			_, err = s.pool.Exec(ctx, `
	delete from carts_items where cart_id=$1 and product_id=$2
`, id, product.ID)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		if product.QTY > count {
			log.Println(ErrNosuch)
			return ErrNosuch
		}
		_, err = s.pool.Exec(ctx, `
update carts_items set  count=count-$3 where product_id=$1 and cart_id=$2
`, product.ID, id, product.QTY)
	}
	return nil
}

func (s *Service) GetCartBYID(ctx context.Context, id int64) ([]*Product, error) {
	items := make([]*Product, 0)
	rows, err := s.pool.Query(ctx, `
select product_id,count from carts_items  where cart_id=$1
`, id)
	for rows.Next() {
		product := &Product{}
		err = rows.Scan(&product.ID, &product.QTY)
		if err != nil {
			return nil, err
		}
		items = append(items, product)
	}
	return items, nil
}

func (s *Service) GetSum(ctx context.Context, products []*Product) (int64, error) {
	sum := int64(0)
	for _, product := range products {
		var cost int64
		err := s.pool.QueryRow(ctx, `
select price from products where id=$1
`, product.ID).Scan(&cost)
		if err != nil {
			return -1, err
		}
		sum += cost
	}
	return sum, nil
}
func (s *Service) DeleteCart(ctx context.Context, ID int64) error {
	_, err := s.pool.Exec(ctx, `
delete from carts_items where cart_id=$1
`, ID)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `
delete from carts where id=$1
`, ID)
	if err != nil {
		return err
	}
	return nil
}
