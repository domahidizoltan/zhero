// Package page is the repository to manage schema pages.
package page

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"

	domain "github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/pkg/collection"
	"github.com/domahidizoltan/zhero/pkg/database"
	"github.com/domahidizoltan/zhero/pkg/paging"
	"github.com/oklog/ulid"
)

const (
	selectPage = `SELECT secondary_identifier, listable_data, data, meta, "references", enabled FROM page WHERE schema_name = ? AND identifier = ?;`
	insertPage = `INSERT INTO page (schema_name, identifier, secondary_identifier, listable_data, data, meta, "references", enabled) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	updatePage = `UPDATE page
		SET secondary_identifier = ?, listable_data = ?, data = ?, meta = ?, "references" = ?, enabled = ?		WHERE schema_name = ? AND identifier = ?;`
	enablePage = `UPDATE page SET enabled = ? WHERE schema_name = ? AND identifier = ?;`
	deletePage = `DELETE FROM page WHERE schema_name = ? AND identifier = ?;`

	insertPageSearch = `INSERT INTO page_search (schema_name, identifier, col0, col1, col2, col3, col4) VALUES (?, ?, ?, ?, ?, ?, ?);`
	updatePageSearch = `UPDATE page_search
		SET col0 = ?, col1 = ?, col2 = ?, col3 = ?, col4 = ?
		WHERE schema_name = ? AND identifier = ?;`
	deletePageSearch = `DELETE FROM page_search WHERE schema_name = ? AND identifier = ?;`

	// selectPageSearch = `SELECT col0,col1,col2,col3,col4 FROM page_search WHERE schema_name = ? AND identifier = ?;`

	listPagesBase  = `SELECT identifier, secondary_identifier, enabled, listable_data, "references" FROM page WHERE schema_name = ?`
	countPagesBase = `SELECT COUNT(*) FROM page WHERE schema_name = ?`

	selectEnabledSchemaNames = `SELECT DISTINCT(schema_name) FROM page WHERE enabled = TRUE ORDER BY schema_name ASC`

	listReferencesByTypeQuery = `SELECT identifier, secondary_identifier, enabled FROM page 
	WHERE schema_name = ? ORDER BY identifier, secondary_identifier`

	listEnabledRoutesByRef = `
SELECT rt.page, rt.route, p.enabled FROM (
  SELECT route, page, ROW_NUMBER() OVER (
    PARTITION BY page ORDER BY version DESC
  ) AS rn FROM route WHERE page IN (?)
) rt 
RIGHT JOIN page p ON rt.page = CONCAT(p.schema_name,'/',p.identifier)
WHERE rt.rn=1`
)

type Repository struct {
	db              *sql.DB
	defaultPageSize uint
}

func NewRepo(db *sql.DB, defaultPageSize uint) *Repository {
	return &Repository{
		db:              db,
		defaultPageSize: defaultPageSize,
	}
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

	metaJSON, err := json.Marshal(page.Meta)
	if err != nil {
		return "", fmt.Errorf("failed to serialize page meta to JSON: %w", err)
	}

	listableDataJSON, err := json.Marshal(page.ListableData)
	if err != nil {
		return "", fmt.Errorf("failed to serialize page listable data to JSON: %w", err)
	}

	referencesJSON, err := json.Marshal(page.References)
	if err != nil {
		return "", fmt.Errorf("failed to serialize page references to JSON: %w", err)
	}

	if _, err := tx.ExecContext(ctx, insertPage,
		page.SchemaName, newID.String(), page.SecondaryIdentifier, listableDataJSON, dataJSON, metaJSON, referencesJSON, page.IsEnabled); err != nil {
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

	metaJSON, err := json.Marshal(page.Meta)
	if err != nil {
		return fmt.Errorf("failed to serialize page meta to JSON: %w", err)
	}

	listableDataJSON, err := json.Marshal(page.ListableData)
	if err != nil {
		return fmt.Errorf("failed to serialize page listable data to JSON: %w", err)
	}

	referencesJSON, err := json.Marshal(page.References)
	if err != nil {
		return fmt.Errorf("failed to serialize page references to JSON: %w", err)
	}

	if _, err := tx.ExecContext(ctx, updatePage,
		page.SecondaryIdentifier, listableDataJSON, dataJSON, metaJSON, referencesJSON, page.IsEnabled, page.SchemaName, identifier); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, updatePageSearch,
		page.SearchVals[0], page.SearchVals[1], page.SearchVals[2], page.SearchVals[3], page.SearchVals[4], page.SchemaName, identifier); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPageBySchemaNameAndIdentifier(ctx context.Context, schemaName, identifier string, onlyEnabled bool) (*domain.Page, error) {
	query := selectPage
	args := []any{schemaName, identifier}

	if onlyEnabled {
		query = strings.Replace(query, ";", " AND enabled = TRUE;", 1)
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	page := domain.Page{
		SchemaName: schemaName,
		Identifier: identifier,
	}
	var dataJSON, metaJSON, listableDataJSON, referencesJSON sql.NullString
	if err := row.Scan(&page.SecondaryIdentifier, &listableDataJSON, &dataJSON, &metaJSON, &referencesJSON, &page.IsEnabled); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(dataJSON.String), &page.Data); err != nil {
		return nil, err
	}

	if metaJSON.Valid && metaJSON.String != "" {
		if err := json.Unmarshal([]byte(metaJSON.String), &page.Meta); err != nil {
			return nil, fmt.Errorf("failed to deserialize page meta: %w", err)
		}
	}

	if listableDataJSON.Valid && listableDataJSON.String != "" {
		if err := json.Unmarshal([]byte(listableDataJSON.String), &page.ListableData); err != nil {
			return nil, fmt.Errorf("failed to deserialize page listable data: %w", err)
		}
	}

	if referencesJSON.Valid && referencesJSON.String != "" {
		if err := json.Unmarshal([]byte(referencesJSON.String), &page.References); err != nil {
			return nil, fmt.Errorf("failed to deserialize page references: %w", err)
		}
	}

	return &page, nil
}

func (r *Repository) List(ctx context.Context, schemaName string, opts domain.ListOptions, onlyEnabled bool) ([]domain.Page, paging.Meta, error) {
	if len(opts.SortBy) == 0 {
		opts.SortBy = "identifier"
	}
	pages := []domain.Page{}

	countQuery := countPagesBase
	countArgs := []any{schemaName}
	query := listPagesBase
	queryArgs := []any{schemaName}

	if onlyEnabled {
		countQuery += " AND enabled = TRUE"
		query += " AND enabled = TRUE"
	}

	if len(opts.SecondaryIdentifierLike) > 0 {
		countQuery += " AND secondary_identifier LIKE ?" // case insensitive in SQlite
		countArgs = append(countArgs, "%"+opts.SecondaryIdentifierLike+"%")
		query += " AND secondary_identifier LIKE ?"
		queryArgs = append(queryArgs, "%"+opts.SecondaryIdentifierLike+"%")
	}
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return pages, paging.Meta{}, fmt.Errorf("failed to count pages: %w", err)
	}

	meta := opts.ToMeta(total, r.defaultPageSize)
	if total == 0 {
		return pages, meta, nil
	}

	query += " ORDER BY " + opts.SortBy + " " + string(opts.SortDir) + " LIMIT ? OFFSET ?"
	queryArgs = append(queryArgs, meta.PageSize, (opts.Page-1)*meta.PageSize)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return pages, meta, fmt.Errorf("failed to query listed pages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Page
		var listableDataJSON, referencesJSON sql.NullString
		if err := rows.Scan(&p.Identifier, &p.SecondaryIdentifier, &p.IsEnabled, &listableDataJSON, &referencesJSON); err != nil {
			return pages, meta, fmt.Errorf("failed to scan listed page row: %w", err)
		}
		if listableDataJSON.Valid && listableDataJSON.String != "" {
			if err := json.Unmarshal([]byte(listableDataJSON.String), &p.ListableData); err != nil {
				return pages, meta, fmt.Errorf("failed to deserialize listable data: %w", err)
			}
		}
		if referencesJSON.Valid && referencesJSON.String != "" {
			if err := json.Unmarshal([]byte(referencesJSON.String), &p.References); err != nil {
				return pages, meta, fmt.Errorf("failed to deserialize page references: %w", err)
			}
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

func (r *Repository) GetPageSchemaNames(ctx context.Context, onlyEnabled bool) ([]string, error) {
	query := selectEnabledSchemaNames
	if !onlyEnabled {
		query = strings.ReplaceAll(query, "WHERE enabled = TRUE ", "")
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for rows.Next() {
		var name string
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil
}

func (r *Repository) ListReferencesByType(ctx context.Context, schemaName string) ([]domain.Page, error) {
	rows, err := r.db.QueryContext(ctx, listReferencesByTypeQuery, schemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []domain.Page{}
	for rows.Next() {
		var ref domain.Page
		if err := rows.Scan(&ref.Identifier, &ref.SecondaryIdentifier, &ref.IsEnabled); err != nil {
			return nil, err
		}
		results = append(results, ref)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *Repository) ListEnabledRoutesByRef(ctx context.Context, refs []string) (map[string]domain.Page, error) {
	if len(refs) == 0 {
		return map[string]domain.Page{}, nil
	}

	query := listEnabledRoutesByRef
	if len(refs) > 1 {
		query = strings.Replace(query, "?", "?"+strings.Repeat(", ?", len(refs)-1), 1)
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	anyRefs := slices.Collect(collection.MapValues(refs, func(in string) any {
		return any(in)
	}))

	rows, err := stmt.Query(anyRefs...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	results := map[string]domain.Page{}
	for rows.Next() {
		var pg domain.Page
		var ref string
		if err := rows.Scan(&ref, &pg.Route, &pg.IsEnabled); err != nil {
			return nil, err
		}
		results[ref] = pg
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
