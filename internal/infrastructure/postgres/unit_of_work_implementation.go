package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhaatha/go-bookshelf/internal/repository"
	"github.com/mhaatha/go-bookshelf/internal/service"
)

// pgxTransaction implements Transaction.
// Transaction is literally a db tx.
type pgxTransaction struct {
	tx pgx.Tx
}

func (t *pgxTransaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *pgxTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

// pgxTransaction creates new GetAuthorRepository and GetBookRepository instance.
// This can be done since pgx.Tx implemented repository.DBTX
func (t *pgxTransaction) GetAuthorRepository() repository.AuthorRepository {
	return repository.NewAuthorRepository(t.tx)
}

func (t *pgxTransaction) GetBookRepository() repository.BookRepository {
	return repository.NewBookRepository(t.tx)
}

// pgxUnitOfWork implements UnitOfWork.
// pgxUnitOfWork is literally a db pool, it holds pgxpool.Pool value inside
// that's why pgxUnitOfWork will be passed in to service parameter.
type pgxUnitOfWork struct {
	db *pgxpool.Pool
}

func NewPgxUnitOfWork(db *pgxpool.Pool) service.UnitOfWork {
	return &pgxUnitOfWork{db: db}
}

func (u *pgxUnitOfWork) Begin(ctx context.Context) (service.Transaction, error) {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		slog.Error("error when creating pgx db transaction", "err", err)
		return nil, err
	}
	return &pgxTransaction{tx: tx}, nil
}
