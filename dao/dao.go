//go:build !solution

package dao

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
)

type dao struct {
	conn *pgx.Conn
}

func (d dao) Create(ctx context.Context, u *User) (UserID, error) {
	var id int
	err := d.conn.QueryRow(
		ctx,
		"INSERT INTO users(name) VALUES($1) RETURNING id",
		u.Name,
	).Scan(&id)
	if err != nil {
		return -1, err
	}
	return UserID(id), nil
}

func (d dao) Update(ctx context.Context, u *User) error {
	_, err := d.conn.Exec(
		ctx,
		"UPDATE users SET name = $1 WHERE id = $2",
		u.Name, u.ID,
	)
	return err
}

func (d dao) Delete(ctx context.Context, id UserID) error {
	_, err := d.conn.Exec(
		ctx,
		"DELETE FROM users WHERE id = $1",
		id,
	)
	return err
}

func (d dao) Lookup(ctx context.Context, id UserID) (User, error) {
	row := d.conn.QueryRow(ctx, "SELECT id, name FROM users WHERE id = $1", id)
	return scanUser(row)
}

func scanUser(row pgx.Row) (User, error) {
	user := User{}
	err := row.Scan(&user.ID, &user.Name)
	// bad way
	if err != nil && err.Error() == "no rows in result set" {
		return user, sql.ErrNoRows
	}
	return user, err
}

func (d dao) List(ctx context.Context) ([]User, error) {
	var res []User

	rows, err := d.conn.Query(ctx, "SELECT id, name FROM users")
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		if user, uErr := scanUser(rows); uErr != nil {
			return res, uErr
		} else {
			res = append(res, user)
		}

	}

	err = rows.Err()
	return res, err
}

func (d dao) Close() error {
	if d.conn != nil {
		ctx := context.Background()
		return d.conn.Close(ctx)
	}
	return nil
}

func CreateDao(ctx context.Context, dsn string) (Dao, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	d := dao{
		conn: conn,
	}

	_, err = conn.Exec(
		ctx,
		`
create table users(
    id SERIAL PRIMARY KEY,
    name varchar(50)
);
`,
	)
	if err != nil {
		_ = conn.Close(ctx)
		return nil, err
	}

	return d, nil
}
