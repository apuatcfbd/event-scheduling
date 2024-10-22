package mysql

import (
	"database/sql"
	"time"
)

func QueryDueEvents(db *sql.DB) (*sql.Rows, error) {
	return db.Query(`
		SELECT id, name, payload 
		FROM jobs 
		WHERE runAt < $1 AND cron = '-'`,
		time.Now(),
	)
}

func QueryScheduleEvent(db *sql.DB, event, payload string, runAt time.Time) (sql.Result, error) {
	return db.Exec(`
		INSERT INTO jobs (name, payload, runAt) VALUES ($1, $2, $3)`,
		event, payload, runAt,
	)
}

func QueryScheduleEventWithCron(db *sql.DB, event, payload, cron string, runAt time.Time) (sql.Result, error) {
	return db.Exec(`
		INSERT INTO jobs (name, payload, runAt, cron) VALUES ($1, $2, $3, $4)`,
		event, payload, runAt, cron,
	)
}

func QueryUpdateEvent(db *sql.DB, event, payload, cron string) (sql.Result, error) {
	return db.Exec(`
		UPDATE jobs SET cron = $1 , payload = $2 
		            WHERE name = $3 AND cron != '-'`,
		cron, payload, event,
	)
}

func QueryDeleteEvent(db *sql.DB, eventId uint) (sql.Result, error) {
	return db.Exec(`DELETE FROM jobs WHERE id = $1`, eventId)
}

func QueryEmptyCronEvents(db *sql.DB) (*sql.Rows, error) {
	return db.Query(`SELECT id, name, payload, cron FROM jobs WHERE cron!='-'`)
}
