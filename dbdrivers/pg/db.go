package pg

import (
	"database/sql"
	"github.com/dipeshdulal/event-scheduling/config"
	_ "github.com/lib/pq"
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
	log.Print("💾 Seeding database with table...")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "public"."jobs" (
			"id"      SERIAL PRIMARY KEY,
			"name"    VARCHAR(50) NOT NULL,
			"payload" text,
			"runAt"   TIMESTAMP NOT NULL,
			"cron"    VARCHAR(50) DEFAULT '-'
		)
	`)

	if err != nil {
		log.Panic("query error: ", err)
	}

	return err
}
