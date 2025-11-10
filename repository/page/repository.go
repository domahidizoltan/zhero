// Package page is the repository to manage schema pages.
package page

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"

	domain "github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/pkg/database"
	"github.com/google/uuid"
)

const (
	maxSearchFields  = 5
	insertPage       = `INSERT INTO page (schema_name, identifier, secondary_identifier, fields, search_columns, enabled) VALUES (?, ?, ?, ?, ?, ?);`
	insertPageSearch = `INSERT INTO page_search (schema_name, identifier, col0, col1, col2, col3, col4) VALUES (?, ?, ?, ?, ?, ?, ?);`
	updatePage       = `UPDATE page 
		SET secondary_identifier = ?, fields = ?, search_columns = ?, enabled = ?
		WHERE schema_name = ? AND identifier = ?;`
	updatePageSearch = `UPDATE page_search 
		SET col0 = ?, col1 = ?, col2 = ?, col3 = ?, col4 = ?
		WHERE schema_name = ? AND identifier = ?;`
	selectPage       = `SELECT secondary_identifier, fields, enabled FROM page WHERE schema_name = ? AND identifier = ?;`
	selectPageSearch = `SELECT col0,col1,col2,col3,col4 FROM page_search WHERE schema_name = ? AND identifier = ?;`
	// searchBySchema     = `SELECT col0,col1,col2,col3,col4 FROM page_search WHERE schema_name = ? MATCH ? ORDER BY rank;`
	// deletePage       = `DELETE FROM page WHERE schema_name = ? AND identifier = ?;`
	// deletePageSearch = `DELETE FROM page_search WHERE schema_name = ? AND identifier = ?;`
)

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// TODO ignore secondaryIdentifier from searchCols because it's searchable as defined
func (r *Repository) Insert(ctx context.Context, page domain.Page) (string, error) {
	tx := database.GetTx(ctx)
	if tx == nil {
		return "", database.ErrTransactionNotFound
	}

	t, err := r.getTransformedFields(page, uuid.NewString())
	if err != nil {
		return "", err
	}

	if _, err := tx.ExecContext(ctx, insertPage,
		page.SchemaName, t.id, t.secID, t.fieldsJSON, t.searchColNamesJSON, page.IsEnabled); err != nil {
		return "", err
	}

	if _, err := tx.ExecContext(ctx, insertPageSearch,
		page.SchemaName, t.id, t.searchVals[0], t.searchVals[1], t.searchVals[2], t.searchVals[3], t.searchVals[4]); err != nil {
		return "", err
	}

	return t.id, nil
}

func (r *Repository) Update(ctx context.Context, identifier string, page domain.Page) error {
	tx := database.GetTx(ctx)
	if tx == nil {
		return database.ErrTransactionNotFound
	}

	t, err := r.getTransformedFields(page, identifier)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, updatePage,
		t.secID, t.fieldsJSON, t.searchColNamesJSON, page.IsEnabled, page.SchemaName, t.id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, updatePageSearch,
		t.searchVals[0], t.searchVals[1], t.searchVals[2], t.searchVals[3], t.searchVals[4], page.SchemaName, t.id); err != nil {
		return err
	}

	return nil
}

type transformedFields struct {
	id, secID, fieldsJSON, searchColNamesJSON string
	searchVals                                []string
}

func (r *Repository) getTransformedFields(page domain.Page, id string) (transformedFields, error) {
	t := transformedFields{
		id: id,
	}

	idx := domain.GetFieldIdxByName(page.Fields, page.Identifier)
	if idx == -1 {
		return t, errors.New("failed to update identifier page field")
	}
	page.Fields[idx].Value = t.id

	if secIdx := domain.GetFieldIdxByName(page.Fields, page.SecondaryIdentifier); secIdx > -1 {
		t.secID = fmt.Sprintf("%s", page.Fields[secIdx].Value)
	}

	fields, searchColNames, searchVals := extractSearchFields(page.Fields)
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return t, err
	}
	searchColNamesJSON, err := json.Marshal(searchColNames)
	if err != nil {
		return t, err
	}

	t.searchVals = searchVals
	t.fieldsJSON = string(fieldsJSON)
	t.searchColNamesJSON = string(searchColNamesJSON)
	return t, nil
}

func extractSearchFields(fields []domain.Field) ([]domain.Field, []string, []string) {
	searchColNames := make([]string, 0, maxSearchFields)
	searchVals := make([]string, 0, maxSearchFields)
	returnFields := make([]domain.Field, 0, len(fields))

	for _, f := range fields {
		if len(f.SearchColumn) > 0 {
			searchColNames = append(searchColNames, f.SearchColumn)
			searchVals = append(searchVals, fmt.Sprintf("%v", f.Value))
			f.Value = ""
		}
		returnFields = append(returnFields, f)
	}

	for range maxSearchFields - len(searchColNames) {
		searchColNames = append(searchColNames, "")
		searchVals = append(searchVals, "")
	}

	return returnFields, searchColNames, searchVals
}

func (r *Repository) GetPageBySchemaNameAndIdentifier(ctx context.Context, schemaName string, identifier string) (*domain.Page, error) {
	row := r.db.QueryRowContext(ctx, selectPage, schemaName, identifier)

	page := domain.Page{
		SchemaName: schemaName,
		Identifier: identifier,
	}
	var fieldsJSON string
	if err := row.Scan(&page.SecondaryIdentifier, &fieldsJSON, &page.IsEnabled); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(fieldsJSON), &page.Fields); err != nil {
		return nil, err
	}

	searchRow := r.db.QueryRowContext(ctx, selectPageSearch, schemaName, identifier)
	if err := searchRow.Err(); err != nil {
		return nil, err
	}

	c := [maxSearchFields]string{}
	if err := searchRow.Scan(&c[0], &c[1], &c[2], &c[3], &c[4]); err != nil {
		return nil, err
	}

	cIdx := 0
	for i, f := range page.Fields {
		if len(f.SearchColumn) > 0 {
			page.Fields[i].Value = c[cIdx]
			cIdx++
		}
	}

	return &page, nil
}
