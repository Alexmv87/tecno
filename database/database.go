package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	uri := os.Getenv("DATABASE_URI")
	return sql.Open("postgres", uri)
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tickets (
			id SERIAL PRIMARY KEY,
			cliente VARCHAR(100) NOT NULL,
			asunto VARCHAR(200) NOT NULL,			
			estado VARCHAR(20) NOT NULL DEFAULT 'abierto',
			creado_en TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Migración completada con éxito")
	return nil
}
