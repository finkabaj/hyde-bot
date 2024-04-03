package postgresql

import (
	"context"
	"fmt"

	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/jackc/pgx/v5"
)

type Postgresql struct {
	db.Database
	Conn *pgx.Conn
}

func (p *Postgresql) setup() (err error) {
	if err = p.Status(); err != nil {
		return
	}

	ctx := context.Background()

	transaction, err := p.Conn.Begin(ctx)

	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			transaction.Rollback(ctx)
		}
	}()

	err = transaction.Commit(ctx)

	return
}

func (p *Postgresql) Connect(credentials *db.DatabaseCredentials) (err error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", credentials.User, credentials.Password, credentials.Host, credentials.Port, credentials.Database)
	p.Conn, err = pgx.Connect(context.Background(), connStr)

	return
}

func (p *Postgresql) Close() {
	p.Conn.Close(context.Background())
}

func (p *Postgresql) Status() (err error) {
	err = p.Conn.Ping(context.Background())

	return
}
