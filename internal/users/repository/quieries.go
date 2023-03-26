package repository

var (
	insertUserQuery = `insert into public.users(id, name, registration_date, balance) overriding user value values (null, :name, now()::timestamp, :balance);`
	deleteUserQuery = `delete from public.users where id = :id;`
	selectUserQuery = `select from public.users where id = :id limit 1;`
)
