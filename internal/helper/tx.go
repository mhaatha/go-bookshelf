package helper

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func CommitOrRollback(ctx context.Context, tx pgx.Tx) {
	if err := tx.Commit(ctx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			log.Fatal(err)
		} else {
			log.Fatal(err)
		}
	}
}
