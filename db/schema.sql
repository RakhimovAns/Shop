CREATE TABLE customers(
    id bigserial primary key ,
    name text not null ,
    phone text not null unique ,
    password text not null ,
    active boolean not null default true,
    balance integer not null,
    created timestamp not null default current_timestamp
);


CREATE TABLE products(
    id bigserial primary key ,
    category text not null ,
    name text not null  unique ,
    price integer not null check ( price>0 ),
    qty integer not null check ( qty>0 ),
    active boolean not null  default true,
    created time not null default  current_timestamp
);


CREATE TABLE purchases(
    id bigserial primary key ,
    customer_id integer references customers,
    product_id integer references products,
    qty integer not null check ( qty>0 ),
    created timestamp not null  default current_timestamp
);


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