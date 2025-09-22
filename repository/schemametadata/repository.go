// Package schemametadata is the repository to manage data blueprint.
package schemametadata

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"

	domain "github.com/domahidizoltan/zhero/domain/schemametadata"
	"github.com/domahidizoltan/zhero/pkg/database"
)

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, schema domain.Schema) (err error) {
	tx := database.GetTx(ctx)
	if tx == nil {
		return database.ErrTransactionNotFound
	}

	if _, err := tx.Exec("DELETE FROM schema_metadata WHERE name = ?", schema.Name); err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO schema_metadata (name, identifier, secondary_identifier) VALUES (?, ?, ?)",
		schema.Name, schema.Identifier, schema.SecondaryIdentifier)
	if err != nil {
		return err
	}

	schemaID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO schema_property_metadata
		(schema_id, name, mandatory, searchable, "type", component, display_order)
		VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, prop := range schema.Properties {
		_, err = stmt.Exec(schemaID, prop.Name, prop.Mandatory, prop.Searchable, prop.Type, prop.Component, prop.Order)
		if err != nil {
			return err
		}
	}

	return err
}
