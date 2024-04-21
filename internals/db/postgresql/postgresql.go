package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgresql struct {
	pool   *pgxpool.Pool
	logger logger.ILogger
}

func NewPostgresql(logger logger.ILogger) *Postgresql {
	return &Postgresql{
		logger: logger,
	}
}

func (p *Postgresql) setup() (err error) {
	if err = p.Status(); err != nil {
		p.logger.Error(err, logger.LogFields{"message": "error while creating tables"})
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
		} else {
			transaction.Commit(ctx)
		}
	}()

	_, err = transaction.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS "users" (
      "userId" VARCHAR(255) PRIMARY KEY,
      "name" VARCHAR(255) NOT NULL
    )
  `)

	if err != nil {
		p.logger.Debug("error creating users table")
		return
	}

	_, err = transaction.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS "guilds" (
      "guildId" VARCHAR(255) PRIMARY KEY,
      "ownerId" VARCHAR(255) NOT NULL
    )
  `)

	if err != nil {
		p.logger.Debug("error creating guilds table")
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
		p.logger.Debug("error creating refreshTokens table")
		return
	}

	_, err = transaction.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS "reactionRules" (
      "emojiId" VARCHAR(255),
      "emojiName" VARCHAR(255),
      "guildId" VARCHAR(255) NOT NULL,
      "ruleAuthor" VARCHAR(255) NOT NULL,
      "actions" INTEGER[] NOT NULL,
      PRIMARY KEY ("guildId", "emojiId", "emojiName"),
      FOREIGN KEY ("guildId") REFERENCES guilds("guildId") ON DELETE CASCADE
    )
  `)

	if err != nil {
		p.logger.Debug("error creating reactionRules table")
		return
	}

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

func (p *Postgresql) CreateGuild(gc guild.GuildCreate) (guild.Guild, error) {
	query := `
    INSERT INTO guilds ("guildId", "ownerId") 
    VALUES ($1, $2) 
    RETURNING *
  `

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	row, err := p.pool.Query(ctx, query, gc.GuildId, gc.OwnerId)
	defer row.Close()

	if err != nil {
		p.logger.Warn(err, logger.LogFields{"message": "error in CreateGuild query"})
		return guild.Guild{}, common.ErrInternal
	}

	newGuild, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[guild.Guild])

	if err != nil {
		p.logger.Error(err, logger.LogFields{"message": "error when collecting rows in CreateGuild"})
		return guild.Guild{}, common.ErrInternal
	}

	return newGuild, nil
}

func (p *Postgresql) GetGuild(guildId string) (guild.Guild, error) {
	query := `
    SELECT * FROM guilds WHERE "guildId" = $1
  `

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	row, err := p.pool.Query(ctx, query, guildId)

	if err != nil {
		p.logger.Error(err, logger.LogFields{"message": "error in GetGuild query"})
		return guild.Guild{}, common.ErrInternal
	}
	defer row.Close()

	foundGuild, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[guild.Guild])

	if err == pgx.ErrNoRows {
		return guild.Guild{}, common.ErrNotFound
	} else if err != nil {
		p.logger.Error(err, logger.LogFields{"message": "error while collecting rows in GetGuild"})
		return guild.Guild{}, common.ErrInternal
	}

	return foundGuild, nil
}

func (p *Postgresql) CreateReactionRules(rules []rule.ReactionRule) ([]rule.ReactionRule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	tx, err := p.pool.Begin(ctx)

	if err != nil {
		p.logger.Error(err, logger.LogFields{"message": "transaction begin in CreateReactionRules"})
		return []rule.ReactionRule{}, common.ErrInternal
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	rows := common.DestructureStructSlice(rules)

	copyCount, err := p.pool.CopyFrom(ctx,
		pgx.Identifier{"reactionRules"},
		[]string{"emojiName", "emojiId", "guildId", "ruleAuthor", "actions"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		p.logger.Error(err, logger.LogFields{"message": "error while inserting to reactionRules"})
		return []rule.ReactionRule{}, common.ErrInternal
	}

	if int(copyCount) != len(rows) {
		err = common.ErrInternal
		p.logger.Error(err, logger.LogFields{"message": fmt.Sprintf("Expected %d but got %d at CreateReactionRules", len(rows), int(copyCount))})
		return []rule.ReactionRule{}, err
	}

	return rules, nil
}

func (p *Postgresql) DeleteReactionRules(ids []string) error {
	return nil
}

func (p *Postgresql) GetReactionRules(gId string) ([]rule.ReactionRule, error) {
	query := `
    SELECT * FROM "reactionRules" WHERE "guildId" = $1
  `

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	rows, err := p.pool.Query(ctx, query, gId)

	if err != nil {
		p.logger.Error(err, logger.LogFields{"message": "error in GetReactionRules query"})
		return []rule.ReactionRule{}, common.ErrInternal
	}

	foundRules, err := pgx.CollectRows(rows, pgx.RowToStructByName[rule.ReactionRule])

	if err == pgx.ErrNoRows {
		return []rule.ReactionRule{}, common.ErrNotFound
	} else if err != nil {
		p.logger.Error(err, logger.LogFields{"message": "error while collecting rows in GetReactionRules"})
		return []rule.ReactionRule{}, common.ErrInternal
	}

	return foundRules, nil
}
