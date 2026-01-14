CREATE TABLE IF NOT EXISTS schema_meta (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    identifier TEXT NOT NULL,
    secondary_identifier TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS schema_meta_properties (
    schema_name TEXT NOT NULL,
    name TEXT NOT NULL,
    mandatory INTEGER DEFAULT 0,
    searchable INTEGER default 0,
    "type" TEXT NOT NULL,
    component TEXT,
    "order" INTEGER NOT NULL DEFAULT 0
);


CREATE TABLE IF NOT EXISTS page (
    schema_name TEXT NOT NULL,
    identifier TEXT NOT NULL,
    secondary_identifier TEXT NOT NULL,
    data TEXT,
    enabled INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS page_schema_id_idx ON page(schema_name, identifier);

CREATE VIRTUAL TABLE IF NOT EXISTS page_search USING FTS5(schema_name, identifier, col0, col1, col2, col3, col4);
