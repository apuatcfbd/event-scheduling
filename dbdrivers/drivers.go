package dbdrivers

import (
	"database/sql"
	"github.com/dipeshdulal/event-scheduling/config"
	"github.com/dipeshdulal/event-scheduling/dbdrivers/mysql"
	"github.com/dipeshdulal/event-scheduling/dbdrivers/pg"
	"log"
)

// GetDbConn returns the db connection for the selected driver (in .env)
func GetDbConn() *sql.DB {
	var conn *sql.DB

	driver := config.EnvDBDriver()
	switch driver {
	case "postgres":
		conn = pg.InitDBConnection()
	case "mysql":
		conn = mysql.InitDBConnection()
	default:
		log.Fatalf("Unsupported DB driver: %s \n", driver)
	}

	return conn
}

// GetDbSeeder returns the seeder fn for the selected db type
func GetDbSeeder() func(db *sql.DB) error {
	driver := config.EnvDBDriver()

	switch driver {
	case "postgres":
		return pg.SeedDB
	case "mysql":
		return mysql.SeedDB
	default:
		log.Fatalf("Unsupported DB driver: %s \n", driver)
	}

	return nil
}
