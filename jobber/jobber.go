package jobber

import (
	"blackbox-v2/internal/notes"
	"blackbox-v2/pkg/logging"
	"time"

	"github.com/go-co-op/gocron"
)

func StartCronServer() {
	t := time.Now()
	logging.CronLogger.Println("Cron server started at ", t)
	s := gocron.NewScheduler(time.UTC)
	s.Every(2).Minute().Do(func() {
		logging.CronLogger.Println("Parsing notes")
		err := notes.ParseNotes()
		if err != nil {
			logging.CronLogger.Println(err)
		}
	})
	s.StartBlocking()
}
