package dbdrivers

import (
	"database/sql"
	"github.com/dipeshdulal/event-scheduling/config"
	"github.com/dipeshdulal/event-scheduling/dbdrivers/mysql"
	"github.com/dipeshdulal/event-scheduling/dbdrivers/pg"
	"time"
)

var dbDriver = config.EnvDBDriver()

// GetDueEventsQuery returns query fn for to get due events from DB
func GetDueEventsQuery() func(db *sql.DB) (*sql.Rows, error) {
	switch dbDriver {
	case "postgres":
		return pg.QueryDueEvents
	case "mysql":
		return mysql.QueryDueEvents
	}

	return nil
}

// GetScheduleEventQuery returns query fn for to add an event in db
func GetScheduleEventQuery() func(db *sql.DB, event, payload string, runAt time.Time) (sql.Result, error) {
	switch dbDriver {
	case "postgres":
		return pg.QueryScheduleEvent
	case "mysql":
		return mysql.QueryScheduleEvent
	}

	return nil
}

// GetScheduleEventWithCronQuery returns query fn for to add an event with cron in db
func GetScheduleEventWithCronQuery() func(db *sql.DB, event, payload, cron string, runAt time.Time) (sql.Result, error) {
	switch dbDriver {
	case "postgres":
		return pg.QueryScheduleEventWithCron
	case "mysql":
		return mysql.QueryScheduleEventWithCron
	}

	return nil
}

// GetUpdateEventQuery returns query fn for to update an event in db
func GetUpdateEventQuery() func(db *sql.DB, event, payload, cron string) (sql.Result, error) {
	switch dbDriver {
	case "postgres":
		return pg.QueryUpdateEvent
	case "mysql":
		return mysql.QueryUpdateEvent
	}

	return nil
}

// GetDeleteEventQuery returns query fn to delete an event from db
func GetDeleteEventQuery() func(db *sql.DB, eventId uint) (sql.Result, error) {
	switch dbDriver {
	case "postgres":
		return pg.QueryDeleteEvent
	case "mysql":
		return mysql.QueryDeleteEvent
	}

	return nil
}

// GetEmptyCronEventsQuery returns query fn to get events that have no cron set
func GetEmptyCronEventsQuery() func(db *sql.DB) (*sql.Rows, error) {
	switch dbDriver {
	case "postgres":
		return pg.QueryEmptyCronEvents
	case "mysql":
		return mysql.QueryEmptyCronEvents
	}

	return nil
}
