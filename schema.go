package gdgberlin

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const schema = `
CREATE EXTENSION "uuid-ossp";

CREATE TYPE water_type AS ENUM (
	'FRESH',
	'SALT',
	'BRACKISH'
);

CREATE TABLE fish (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),

	name TEXT NOT NULL,
	fin_count INT NOT NULL,
	water_type water_type NOT NULL DEFAULT 'SALT'::water_type
);
`

func NewDB(pgDSN string) (*sql.DB, error) {
	db, err := sql.Open("postgres", pgDSN)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
