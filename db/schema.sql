DROP table purchase;
CREATE TABLE purchases(
    id bigserial primary key ,
    customer_id integer references customers,
    product_id integer references products,
    qty integer not null check ( qty>0 ),
    created timestamp not null  default current_timestamp
);