package repository

var (
	updateUserQuery        = `update public.users set balance = balance + :balance_diff where id = :user_id returning *;`
	insertTransactionQuery = `insert into public.transactions(id, user_id, time, balance_diff) overriding user value values (null, :user_id, now()::timestamp, :balance_diff);`
)
