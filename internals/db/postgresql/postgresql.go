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

func (p *Postgresql) Connect(credentials *db.DatabaseCredentials) error {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", credentials.User, credentials.Password, credentials.Host, credentials.Port, credentials.Database)
	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {
		return err
	}

	p.Conn = conn

	return nil
}

func (p *Postgresql) Close() {
	p.Conn.Close(context.Background())
}

func (p *Postgresql) Status() error {
	if err := p.Conn.Ping(context.Background()); err != nil {
		return err
	}

	return nil
}
