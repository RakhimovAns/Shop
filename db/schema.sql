CREATE TABLE cart(
                     id bigserial primary key ,
                     customer_id bigserial references customers(id),
                     items product_qty[],
                     created timestamp not null  default current_timestamp
)