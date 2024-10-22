package main

import (
	"context"
	"database/sql"
	"github.com/dipeshdulal/event-scheduling/dbdrivers"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler data structure
type Scheduler struct {
	db          *sql.DB
	listeners   Listeners
	cron        *cron.Cron
	cronEntries map[string]cron.EntryID
}

// Listeners has attached event listeners
type Listeners map[string]ListenFunc

// ListenFunc function that listens to events
type ListenFunc func(string)

// Event structure
type Event struct {
	ID      uint
	Name    string
	Payload string
	Cron    string
}

// NewScheduler creates a new scheduler
func NewScheduler(db *sql.DB, listeners Listeners) Scheduler {

	return Scheduler{
		db:          db,
		listeners:   listeners,
		cron:        cron.New(),
		cronEntries: map[string]cron.EntryID{},
	}

}

// AddListener adds the listener function to Listeners
func (s Scheduler) AddListener(event string, listenFunc ListenFunc) {
	s.listeners[event] = listenFunc
}

// callListeners calls the event listener of provided event
func (s Scheduler) callListeners(event Event) {
	eventFn, ok := s.listeners[event.Name]
	if ok {
		go eventFn(event.Payload)
		_, err := dbdrivers.GetDeleteEventQuery()(s.db, event.ID)
		if err != nil {
			log.Print("💀 error: ", err)
		}
	} else {
		log.Print("💀 error: couldn't find event listeners attached to ", event.Name)
	}

}

// CheckEventsInInterval checks the event in given interval
func (s Scheduler) CheckEventsInInterval(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("⏰ Ticks Received...")
				events := s.checkDueEvents()
				for _, e := range events {
					s.callListeners(e)
				}
			}

		}
	}()
}

// checkDueEvents checks and returns due events
func (s Scheduler) checkDueEvents() []Event {
	events := []Event{}
	rows, err := dbdrivers.GetDueEventsQuery()(s.db)
	if err != nil {
		log.Print("💀 error: ", err)
		return nil
	}
	for rows.Next() {
		evt := Event{}
		rows.Scan(&evt.ID, &evt.Name, &evt.Payload)
		events = append(events, evt)
	}
	return events
}

// Schedule schedules the provided events
func (s Scheduler) Schedule(event string, payload string, runAt time.Time) {
	log.Print("🚀 Scheduling event ", event, " to run at ", runAt)
	_, err := dbdrivers.GetScheduleEventQuery()(s.db, event, payload, runAt)
	if err != nil {
		log.Print("schedule insert error: ", err)
	}
}

// ScheduleCron schedules a cron job
func (s Scheduler) ScheduleCron(event string, payload string, cron string) {
	log.Print("🚀 Scheduling event ", event, " with cron string ", cron)
	entryID, ok := s.cronEntries[event]
	if ok {
		s.cron.Remove(entryID)
		_, err := dbdrivers.GetUpdateEventQuery()(s.db, cron, payload, event)
		if err != nil {
			log.Print("schedule cron update error: ", err)
		}
	} else {
		_, err := dbdrivers.GetScheduleEventWithCronQuery()(s.db, event, payload, cron, time.Now())
		if err != nil {
			log.Print("schedule cron insert error: ", err)
		}
	}

	eventFn, ok := s.listeners[event]
	if ok {
		entryID, err := s.cron.AddFunc(cron, func() { eventFn(payload) })
		s.cronEntries[event] = entryID
		if err != nil {
			log.Print("💀 error: ", err)
		}
	}
}

// attachCronJobs attaches cron jobs
func (s Scheduler) attachCronJobs() {
	log.Printf("Attaching cron jobs")
	rows, err := dbdrivers.GetEmptyCronEventsQuery()(s.db)
	if err != nil {
		log.Print("💀 error: ", err)
	}
	for rows.Next() {
		evt := Event{}
		rows.Scan(&evt.ID, &evt.Name, &evt.Payload, &evt.Cron)
		eventFn, ok := s.listeners[evt.Name]
		if ok {
			entryID, err := s.cron.AddFunc(evt.Cron, func() { eventFn(evt.Payload) })
			s.cronEntries[evt.Name] = entryID

			if err != nil {
				log.Print("💀 error: ", err)
			}
		}
	}
}

// StartCron starts cron job
func (s Scheduler) StartCron() func() {
	s.attachCronJobs()
	s.cron.Start()

	return func() {
		s.cron.Stop()
	}
}
