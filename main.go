package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	"github.com/dipeshdulal/event-scheduling/customevents"
)

var eventListeners = Listeners{
	"SendEmail": customevents.SendEmail,
	"PayBills":  customevents.PayBills,
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file ", err)
	}

	db := initDBConnection()
	seedDB(db)

	scheduler := NewScheduler(db, eventListeners)
	scheduler.CheckDueEvents()

	scheduler.Schedule("SendEmail", "mail: nilkantha.dipesh@gmail.com", time.Now().Add(1*time.Minute))

	go func() {
		for range interrupt {
			log.Println("\n❌ Interrupt received closing...")
			cancel()
		}
	}()

	<-ctx.Done()
}
