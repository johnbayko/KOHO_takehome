-- PRAGMA foreign_keys = ON;

drop table if exists transactions;
drop table if exists accounts;
drop table if exists customers;


create table customers (
    customer_id text primary key not null
);

create table accounts (
    customer_id text primary key references customers(customer_id) not null,
    balance integer not null
);

create table transactions (
    id text primary key not null,
    customer_id text references customers(customer_id) not null,
    load_amount integer not null,
    time text not null
);

