package main

import "database/sql"

func databaseSchema(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}

	// server_init.sql
	sqlQuery := `CREATE TABLE IF NOT EXISTS public.user
(
    id       BIGSERIAL PRIMARY KEY,
    login varchar,
    password varchar,
    balance_current NUMERIC(7,2) DEFAULT 0 NOT NULL,
    balance_withdraw NUMERIC(7,2) DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS public.order
(
    id  		BIGINT PRIMARY KEY,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id 	BIGINT,
    status 		SMALLINT,
    sum NUMERIC(7,2),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id)
);

CREATE TABLE IF NOT EXISTS public.withdrawal
(
    id  BIGSERIAL PRIMARY KEY,
    processed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    order_id BIGINT,
    user_id BIGINT,
    sum NUMERIC(7,2),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id)
);`
	_, err := db.Exec(sqlQuery)
	return err
}
