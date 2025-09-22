package schemametadata

import (
	"database/sql"

	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"

	sqliteDdl "github.com/domahidizoltan/zhero/data/db/sqlite"
	domain "github.com/domahidizoltan/zhero/service/schemametadata"
)

// Schema represents a schema definition to be saved
type Repository struct {
	db *sql.DB
}

func New(dbFile string) (*Repository, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	// TODO extract migration
	if _, err := db.Exec(sqliteDdl.SchemametadataDdl); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() {
	if err := r.db.Close(); err != nil {
		log.Err(err).Msg("failed to close database")
	}
}

func (r *Repository) SaveSchema(schema domain.Schema) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Delete existing schema and properties if it exists, to handle updates cleanly
	_, err = tx.Exec("DELETE FROM schema_metadata WHERE name = ?", schema.Name)
	if err != nil {
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
