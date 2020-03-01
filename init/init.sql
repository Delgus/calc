CREATE TABLE records (
    id serial not null,
    name varchar(255) not null,
    numeric_param numeric(20,6),
    real_param real
);