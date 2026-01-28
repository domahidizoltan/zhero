// Package schema is the repository to manage data blueprint.
package schema

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	domain "github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/database"
)

const (
	upsertSchemaMeta = `
		INSERT INTO schema_meta (name, identifier, secondary_identifier)
		VALUES (?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			identifier = excluded.identifier,
			secondary_identifier = excluded.secondary_identifier;
	`
	selectSchemaMetaByName = `SELECT name, identifier, secondary_identifier FROM schema_meta WHERE name = ?;`
	selectSchemaMetaNames  = `SELECT name FROM schema_meta ORDER BY name asc;`

	deleteSchemaMetaProps             = `DELETE FROM schema_meta_properties WHERE schema_name = ?;`
	insertSchemaMetaPropsPrefix       = `INSERT INTO schema_meta_properties (schema_name, name, mandatory, searchable, [type], component, [order]) VALUES `
	selectSchemaMetaPropsBySchemaName = `
		SELECT name, mandatory, searchable, [type], component, [order]
		FROM schema_meta_properties
		WHERE schema_name = ?
		ORDER BY [order] ASC;
	`
)

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(ctx context.Context, schema domain.SchemaMeta) error {
	tx := database.GetTx(ctx)
	if tx == nil {
		return database.ErrTransactionNotFound
	}

	if _, err := tx.ExecContext(ctx, upsertSchemaMeta, schema.Name, schema.Identifier, schema.SecondaryIdentifier); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, deleteSchemaMetaProps, schema.Name); err != nil {
		return err
	}

	insertProps := insertSchemaMetaPropsPrefix + strings.Repeat("(?, ?, ?, ?, ?, ?, ?),", len(schema.Properties))
	insertProps = insertProps[:len(insertProps)-1] + ";"
	propValues := make([]any, 0, len(schema.Properties)*7)
	for _, prop := range schema.Properties {
		propValues = append(propValues, schema.Name, prop.Name, prop.Mandatory, prop.Searchable, prop.Type, prop.Component, prop.Order)
	}

	if _, err := tx.ExecContext(ctx, insertProps, propValues...); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetByClassName(ctx context.Context, name string) (*domain.SchemaMeta, error) {
	row := r.db.QueryRowContext(ctx, selectSchemaMetaByName, name)

	var schema domain.SchemaMeta
	if err := row.Scan(&schema.Name, &schema.Identifier, &schema.SecondaryIdentifier); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, selectSchemaMetaPropsBySchemaName, name)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}

		var prop domain.Property
		if err := rows.Scan(&prop.Name, &prop.Mandatory, &prop.Searchable, &prop.Type, &prop.Component, &prop.Order); err != nil {
			return nil, err
		}
		schema.Properties = append(schema.Properties, prop)
	}

	return &schema, nil
}

func (r *Repository) GetAllNames(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, selectSchemaMetaNames)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return names, nil
}
