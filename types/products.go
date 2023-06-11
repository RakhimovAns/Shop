package types

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
