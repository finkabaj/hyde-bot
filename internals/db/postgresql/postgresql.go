package postgresql

import (
	"context"
	"fmt"

	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/finkabaj/hyde-bot/internals/logger"
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

	_, err = transaction.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS users (
      user_id VARCHAR(255) PRIMARY KEY,
      name VARCHAR(255) NOT NULL
    )
  `)

	if err != nil {
		logger.Debug("error creating users table")
		return
	}

	_, err = transaction.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS guilds (
      guild_id VARCHAR(255) PRIMARY KEY,
      name VARCHAR(255) NOT NULL
    )
  `)

	if err != nil {
		logger.Debug("error creating guilds table")
		return
	}

	_, err = transaction.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS refresh_tokens (
      user_id VARCHAR(255) PRIMARY KEY,
      token TEXT NOT NULL,
      expires DATE NOT NULL,
      CONSTRAINT fk_user
        FOREIGN KEY(user_id)
          REFERENCES users(user_id) ON DELETE CASCADE
    )
  `)

	if err != nil {
		logger.Debug("error creating refresh_tokens table")
		return
	}

	err = transaction.Commit(ctx)

	return
}

func (p *Postgresql) Connect(credentials *db.DatabaseCredentials) (err error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", credentials.User, credentials.Password, credentials.Host, credentials.Port, credentials.Database)
	p.Conn, err = pgx.Connect(context.Background(), connStr)

	if err != nil {
		return
	}

	err = p.setup()

	return
}

func (p *Postgresql) Close() {
	p.Conn.Close(context.Background())
}

func (p *Postgresql) Status() (err error) {
	err = p.Conn.Ping(context.Background())

	return
}
