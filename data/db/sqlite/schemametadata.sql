-- schema_metadata stores the high-level information about a defined schema.
CREATE TABLE IF NOT EXISTS schema_metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    identifier TEXT,
    secondary_identifier TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- schema_property_metadata stores properties for each schema defined in schema_metadata.
CREATE TABLE IF NOT EXISTS schema_property_metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    schema_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    mandatory BOOLEAN NOT NULL,
    searchable BOOLEAN NOT NULL,
    "type" TEXT NOT NULL,
    component TEXT NOT NULL,
    display_order INTEGER NOT NULL,
    FOREIGN KEY (schema_id) REFERENCES schema_metadata(id) ON DELETE CASCADE
);
