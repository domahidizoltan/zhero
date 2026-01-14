// Package page is the repository to manage schema pages.
package page

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	domain "github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/pkg/database"
	"github.com/oklog/ulid"
)

const (
	insertPage       = `INSERT INTO page (schema_name, identifier, secondary_identifier, data, enabled) VALUES (?, ?, ?, ?, ?);`
	insertPageSearch = `INSERT INTO page_search (schema_name, identifier, col0, col1, col2, col3, col4) VALUES (?, ?, ?, ?, ?, ?, ?);`
	updatePage       = `UPDATE page 
		SET secondary_identifier = ?, data = ?, enabled = ?
		WHERE schema_name = ? AND identifier = ?;`
	updatePageSearch = `UPDATE page_search 
		SET col0 = ?, col1 = ?, col2 = ?, col3 = ?, col4 = ?
		WHERE schema_name = ? AND identifier = ?;`
	selectPage = `SELECT secondary_identifier, data, enabled FROM page WHERE schema_name = ? AND identifier = ?;`
	// selectPageSearch = `SELECT col0,col1,col2,col3,col4 FROM page_search WHERE schema_name = ? AND identifier = ?;`

	listPages  = `SELECT identifier, secondary_identifier, enabled FROM page WHERE schema_name = ?`
	countPages = `SELECT COUNT(*) FROM page WHERE schema_name = ?`

	enablePage       = `UPDATE page SET enabled = ? WHERE schema_name = ? AND identifier = ?;`
	deletePage       = `DELETE FROM page WHERE schema_name = ? AND identifier = ?;`
	deletePageSearch = `DELETE FROM page_search WHERE schema_name = ? AND identifier = ?;`
)

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(ctx context.Context, page domain.Page, idField string) (string, error) {
	tx := database.GetTx(ctx)
	if tx == nil {
		return "", database.ErrTransactionNotFound
	}

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	newID, err := ulid.New(ms, entropy)
	if err != nil {
		return "", err
	}

	page.Data[idField] = newID.String()
	page.Data["@id"] = newID.String()
	page.Data["@type"] = page.SchemaName
	dataJSON, err := json.Marshal(page.Data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize page data to JSON: %w", err)
	}

	if _, err := tx.ExecContext(ctx, insertPage,
		page.SchemaName, newID.String(), page.SecondaryIdentifier, dataJSON, page.IsEnabled); err != nil {
		return "", err
	}

	if _, err := tx.ExecContext(ctx, insertPageSearch,
		page.SchemaName, newID.String(), page.SearchVals[0], page.SearchVals[1], page.SearchVals[2], page.SearchVals[3], page.SearchVals[4]); err != nil {
		return "", err
	}

	return newID.String(), nil
}

func (r *Repository) Update(ctx context.Context, identifier string, page domain.Page, idField string) error {
	tx := database.GetTx(ctx)
	if tx == nil {
		return database.ErrTransactionNotFound
	}

	page.Data[idField] = identifier
	page.Data["@id"] = identifier
	page.Data["@type"] = page.SchemaName
	dataJSON, err := json.Marshal(page.Data)
	if err != nil {
		return fmt.Errorf("failed to serialize page data to JSON: %w", err)
	}

	if _, err := tx.ExecContext(ctx, updatePage,
		page.SecondaryIdentifier, dataJSON, page.IsEnabled, page.SchemaName, identifier); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, updatePageSearch,
		page.SearchVals[0], page.SearchVals[1], page.SearchVals[2], page.SearchVals[3], page.SearchVals[4], page.SchemaName, identifier); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPageBySchemaNameAndIdentifier(ctx context.Context, schemaName, identifier string) (*domain.Page, error) {
	row := r.db.QueryRowContext(ctx, selectPage, schemaName, identifier)

	page := domain.Page{
		SchemaName: schemaName,
		Identifier: identifier,
	}
	var dataJSON string
	if err := row.Scan(&page.SecondaryIdentifier, &dataJSON, &page.IsEnabled); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(dataJSON), &page.Data); err != nil {
		return nil, err
	}

	return &page, nil
}

const defaultPageSize uint = 20

func (r *Repository) List(ctx context.Context, schemaName string, opts domain.ListOptions) ([]domain.Page, domain.PagingMeta, error) {
	pageSize := defaultPageSize
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize > 0 {
		pageSize = opts.PageSize
	}
	if len(opts.SortBy) == 0 {
		opts.SortBy = "identifier"
	}
	meta := domain.PagingMeta{
		PageSize:    uint(pageSize),
		CurrentPage: opts.Page,
	}
	pages := []domain.Page{}

	countQuery := countPages
	countArgs := []any{schemaName}
	query := listPages
	queryArgs := []any{schemaName}

	if len(opts.SecondaryIdentifierLike) > 0 {
		countQuery += " AND secondary_identifier LIKE ?" // case insensitive in SQlite
		countArgs = append(countArgs, "%"+opts.SecondaryIdentifierLike+"%")
		query += " AND secondary_identifier LIKE ?"
		queryArgs = append(queryArgs, "%"+opts.SecondaryIdentifierLike+"%")
	}
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return pages, meta, fmt.Errorf("failed to count pages: %w", err)
	}
	if total == 0 {
		return pages, meta, nil
	}

	meta.TotalItems = uint(total)
	meta.TotalPages = uint(meta.TotalItems / meta.PageSize)

	query += " ORDER BY " + opts.SortBy + " " + string(opts.SortDir) + " LIMIT ? OFFSET ?"
	queryArgs = append(queryArgs, pageSize, (opts.Page-1)*pageSize)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return pages, meta, fmt.Errorf("failed to query listed pages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Page
		if err := rows.Scan(&p.Identifier, &p.SecondaryIdentifier, &p.IsEnabled); err != nil {
			return pages, meta, fmt.Errorf("failed to scan listed page row: %w", err)
		}
		pages = append(pages, p)
	}

	if err = rows.Err(); err != nil {
		return pages, meta, fmt.Errorf("error during rows iteration: %w", err)
	}

	return pages, meta, nil
}

func (r *Repository) Enable(ctx context.Context, schemaName, identifier string, enable bool) error {
	tx := database.GetTx(ctx)
	if tx == nil {
		return database.ErrTransactionNotFound
	}

	_, err := tx.ExecContext(ctx, enablePage, enable, schemaName, identifier)
	return err
}

func (r *Repository) Delete(ctx context.Context, schemaName, identifier string) error {
	tx := database.GetTx(ctx)
	if tx == nil {
		return database.ErrTransactionNotFound
	}

	if _, err := tx.ExecContext(ctx, deletePageSearch, schemaName, identifier); err != nil {
		return err
	}

	_, err := tx.ExecContext(ctx, deletePage, schemaName, identifier)
	return err
}
