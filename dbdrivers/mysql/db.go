package mysql

import (
	"database/sql"
	"github.com/dipeshdulal/event-scheduling/config"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func InitDBConnection() *sql.DB {
	driver, connStr := config.EnvDBDriver(), config.EnvDBDsn()
	db, err := sql.Open(driver, connStr)

	if err != nil {
		log.Panic("couldn't connect to database", err)
	}

	return db
}

func SeedDB(db *sql.DB) error {
	log.Printf("💾 [%s] Seeding database with table...", config.EnvDBDriver())
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS jobs (
			id      SERIAL PRIMARY KEY,
			name    VARCHAR(50) NOT NULL,
			payload text,
			runAt   TIMESTAMP NOT NULL,
			cron    VARCHAR(50) DEFAULT '-'
		)
	`)

	if err != nil {
		log.Panic("query error: ", err)
	}

	return err
}
