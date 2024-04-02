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
