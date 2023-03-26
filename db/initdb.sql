create user postgres superuser;

create table public.users(
    id int generated always as identity primary key,
    name text unique not null,
    registration_date timestamp not null,
    balance int check (balance >= 0)
);

create table public.transactions(
    id int generated always as identity primary key,
    user_id int,
    foreign key (user_id) references public.users(id),
    time timestamp not null,
    balance_diff int
);
