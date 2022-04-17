//go:build !solution

package ledger

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type ledger struct {
	db *sql.DB
}

func (l ledger) CreateAccount(ctx context.Context, id ID) error {
	_, err := l.db.ExecContext(
		ctx,
		"INSERT INTO accounts(id) VALUES($1)",
		id,
	)
	return err
}

func (l ledger) GetBalance(ctx context.Context, id ID) (Money, error) {
	var balance Money
	err := l.db.QueryRowContext(
		ctx,
		"SELECT balance FROM accounts WHERE id = $1",
		id,
	).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

type exec interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func deposit(db exec, ctx context.Context, id ID, amount Money) error {
	_, err := db.ExecContext(
		ctx,
		"UPDATE accounts SET balance = (balance + $1) WHERE id = $2",
		amount,
		id,
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.ConstraintName == "accounts_balance_check" {
			return ErrNoMoney
		}
	}
	return err
}

func (l ledger) Deposit(ctx context.Context, id ID, amount Money) error {
	return deposit(l.db, ctx, id, amount)
}

func (l ledger) Withdraw(ctx context.Context, id ID, amount Money) error {
	return l.Deposit(ctx, id, -amount)
}

func (l ledger) Transfer(ctx context.Context, from, to ID, amount Money) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	rows, err := tx.QueryContext(
		ctx,
		"SELECT balance FROM accounts WHERE id IN ($1, $2) ORDER BY id FOR UPDATE SKIP LOCKED",
		from,
		to,
	)
	if err != nil {
		return ErrNoMoney
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		cnt++
		_ = rows.Scan()
	}
	if rows.Err() != nil {
		return ErrNoMoney
	}

	if cnt < 2 {
		return ErrNoMoney
	}

	//	getBal := func(id ID) (Money, error) {
	//		var balance Money
	//		err = l.db.QueryRowContext(
	//			ctx,
	//			`
	//SELECT balance FROM accounts WHERE id = $1
	//FOR UPDATE SKIP LOCKED
	//LIMIT 1
	//`,
	//			id,
	//		).Scan(&balance)
	//		return balance, err
	//	}
	//
	//	fromBal, err := getBal(from)
	//	if err != nil {
	//		return err
	//	}
	//	if fromBal < amount {
	//		return ErrNoMoney
	//	}
	err = deposit(tx, ctx, from, -amount)
	if err != nil {
		return err
	}
	err = deposit(tx, ctx, to, amount)
	if err != nil {
		return err
	}

	//_, err = tx.ExecContext(ctx,
	//	`call transfer($1, $2, $3);`,
	//	from, to,
	//	amount,
	//)
	//var pgErr *pgconn.PgError
	//if errors.As(err, &pgErr) {
	//	if pgErr.ConstraintName == "accounts_balance_check" {
	//		return ErrNoMoney
	//	}
	//}
	//if err != nil {
	//	return err
	//}

	return tx.Commit()
}

func (l ledger) Close() error {
	if l.db != nil {
		return l.db.Close()
	}
	return nil
}

func New(ctx context.Context, dsn string) (Ledger, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	res := ledger{
		db: db,
	}
	_, err = db.ExecContext(
		ctx,
		`
create table accounts(
    id TEXT PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0 CHECK(balance >= 0)
);
create or replace procedure transfer(
   sender TEXT,
   receiver TEXT, 
   amount BIGINT
)
language plpgsql    
as $$
begin
    -- subtracting the amount from the sender's account 
    update accounts 
    set balance = balance - amount 
    where id = sender;

    -- adding the amount to the receiver's account
    update accounts 
    set balance = balance + amount 
    where id = receiver;
end;$$

`)
	if err != nil {
		defer db.Close()
		return nil, err
	}
	return res, nil
}
