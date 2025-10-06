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
