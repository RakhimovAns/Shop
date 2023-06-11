package types

import "time"

type Purchase struct {
	ID         int64     `json:"id"`
	CustomerId int64     `json:"customer_id"`
	ProductId  int64     `json:"product_id"`
	QTY        int64     `json:"qty"`
	Created    time.Time `json:"created"`
}
