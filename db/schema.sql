drop table carts_items,carts;
CREATE TABLE carts
(
  id         BIGSERIAL PRIMARY KEY,
  customer_id    integer references customers (id) unique ,
  created_at timestamp not null DEFAULT current_timestamp
);

CREATE TABLE carts_items
(
  cart_id INTEGER REFERENCES carts (id)  ,
  product_id integer REFERENCES products (id) unique ,
  count   INTEGER NOT NULL,
      CONSTRAINT carts_items_pkey PRIMARY KEY (cart_id, product_id)
);

