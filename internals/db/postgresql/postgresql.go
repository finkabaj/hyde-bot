package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgresql struct {
	pool *pgxpool.Pool
}

func (p *Postgresql) setup() (err error) {
	if err = p.Status(); err != nil {
		return
	}

	ctx := context.Background()

	transaction, err := p.pool.Begin(ctx)

	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			transaction.Rollback(ctx)
		}
	}()

	_, err = transaction.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS "users" (
      "userId" VARCHAR(255) PRIMARY KEY,
      "name" VARCHAR(255) NOT NULL
    )
  `)

	if err != nil {
		logger.Debug("error creating users table")
		return
	}

	_, err = transaction.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS "guilds" (
      "guildId" VARCHAR(255) PRIMARY KEY,
      "ownerId" VARCHAR(255) NOT NULL
    )
  `)

	if err != nil {
		logger.Debug("error creating guilds table")
		return
	}

	_, err = transaction.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS "refreshTokens" (
      "userId" VARCHAR(255) PRIMARY KEY,
      "token" TEXT NOT NULL,
      "expires" DATE NOT NULL,
      CONSTRAINT "fkUser"
        FOREIGN KEY("userId")
          REFERENCES users("userId") ON DELETE CASCADE
    )
  `)

	if err != nil {
		logger.Debug("error creating refreshTokens table")
		return
	}

	err = transaction.Commit(ctx)

	return
}

func (p *Postgresql) Connect(credentials db.DatabaseCredentials) (err error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", credentials.User, credentials.Password, credentials.Host, credentials.Port, credentials.Database)
	p.pool, err = pgxpool.New(context.Background(), connStr)

	if err != nil {
		return
	}

	err = p.setup()

	return
}

func (p *Postgresql) Close() {
	p.pool.Close()
}

func (p *Postgresql) Status() (err error) {
	err = p.pool.Ping(context.Background())

	return
}

func (p *Postgresql) CreateGuild(guild *guild.GuildCreate) (*guild.Guild, error) {
	return nil, nil
}

func (p *Postgresql) GetGuild(guildId string) (*guild.Guild, error) {
	query := `
    SELECT * FROM guilds WHERE "guildId" = $1
  `

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	row, err := p.pool.Query(ctx, query, guildId)
	defer row.Close()

	if err != nil {
		return nil, err
	}

	foundGuild, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[guild.Guild])

	if err != nil && err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &foundGuild, nil
}
