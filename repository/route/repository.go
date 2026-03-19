// Package route is the repository to manage custom URL routes.
package route

import (
	"context"
	"database/sql"
	"errors"

	domain "github.com/domahidizoltan/zhero/domain/route"
	"github.com/domahidizoltan/zhero/pkg/database"
)

const (
	selectLatestRouteByRoute  = `SELECT route, page, version FROM route WHERE route = ?;`
	selectLatestVersionByPage = `SELECT route, page, version FROM route WHERE page = ? ORDER BY version DESC LIMIT 1;`
	insertRoute               = `INSERT INTO route (route, page, version) VALUES (?, ?, (SELECT COALESCE(MAX(version), 0) + 1 FROM route WHERE page = ?));`
)

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, route, page string) error {
	tx := database.GetTx(ctx)
	if tx == nil {
		return database.ErrTransactionNotFound
	}

	_, err := tx.ExecContext(ctx, insertRoute, route, page, page)
	return err
}

func (r *Repository) GetByRoute(ctx context.Context, route string) (*domain.Route, error) {
	row := r.db.QueryRowContext(ctx, selectLatestRouteByRoute, route)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var rt domain.Route
	if err := row.Scan(&rt.Route, &rt.Page, &rt.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &rt, nil
}

func (r *Repository) GetLatestVersion(ctx context.Context, page string) (*domain.Route, error) {
	row := r.db.QueryRowContext(ctx, selectLatestVersionByPage, page)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var rt domain.Route
	if err := row.Scan(&rt.Route, &rt.Page, &rt.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &rt, nil
}
