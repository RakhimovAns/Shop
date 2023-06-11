package postgresql

import (
	"context"
	"github.com/RakhimovAns/Shop/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type ProductService struct {
	pool *pgxpool.Pool
}

func NewProductService(pool *pgxpool.Pool) *ProductService {
	return &ProductService{pool: pool}
}

func (s *ProductService) AllActiveProducts(ctx context.Context) ([]*types.Product, error) {
	items := make([]*types.Product, 0)
	rows, err := s.pool.Query(ctx, `
		select id,category,name,price,qty from products where  active=true order by id limit  500 
`)
	if err != nil {

		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		item := &types.Product{}
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

func (s *ProductService) AllCategories(ctx context.Context) ([]*string, error) {
	Categories := make([]*string, 0)
	rows, err := s.pool.Query(ctx, `
		select distinct category from products where active=true limit 500 
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
func (s *ProductService) GetByCategory(ctx context.Context, category string) ([]*types.Product, error) {
	items := make([]*types.Product, 0)
	rows, err := s.pool.Query(ctx, `
		select id,category,name,price,qty from products where active=true and category=$1 limit 500
`, category)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		product := &types.Product{}
		err = rows.Scan(&product.ID, &product.Category, &product.Name, &product.Price, &product.QTY)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, product)
	}
	return items, nil
}

func (s *ProductService) Search(ctx context.Context, Key string) ([]*types.Product, error) {
	items := make([]*types.Product, 0)
	rows, err := s.pool.Query(ctx, `
		select id,category,name,price,qty from products where  active=true and lower(name) like lower($1)
`, Key)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		product := &types.Product{}
		err = rows.Scan(&product.ID, &product.Category, &product.Name, &product.Price, &product.QTY)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, product)
	}
	return items, nil
}
