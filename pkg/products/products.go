package product

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
	ID       int64  `json:"id"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	QTY      int64  `json:"qty"`
}
type ProductsWithoutCategory struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	QTY   int64  `json:"qty"`
}

func (s *Service) AllActiveProducts(ctx context.Context) ([]*Product, error) {
	items := make([]*Product, 0)
	rows, err := s.pool.Query(ctx, `
		select id,category,name,price,qty from products where  active=true order by category limit  500 
`)
	if err != nil {

		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		item := &Product{}
		err = rows.Scan(&item.ID, &item.Category, &item.Name, &item.Price, &item.QTY)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return items, nil
}

func (s *Service) AllCategories(ctx context.Context) ([]*string, error) {
	Categories := make([]*string, 0)
	rows, err := s.pool.Query(ctx, `
		select category from products where active=true limit 500
`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var category *string
		err = rows.Scan(&category)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		Categories = append(Categories, category)
	}
	return Categories, nil
}
func (s *Service) GetByCategory(ctx context.Context, category string) ([]*ProductsWithoutCategory, error) {
	items := make([]*ProductsWithoutCategory, 0)
	rows, err := s.pool.Query(ctx, `
		select id,name,price,qty from products where active=true and category=$1 limit 500
`, category)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		product := &ProductsWithoutCategory{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.QTY)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, product)
	}
	return items, nil
}

func (s *Service) Search(ctx context.Context, Key string) ([]*Product, error) {
	items := make([]*Product, 0)
	rows, err := s.pool.Query(ctx, `
		select id,category,name,price,qty from products where  active=true and lower(name) like lower($1)
`, Key)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		product := &Product{}
		err = rows.Scan(&product.ID, &product.Category, &product.Name, &product.Price, &product.QTY)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, product)
	}
	return items, nil
}
