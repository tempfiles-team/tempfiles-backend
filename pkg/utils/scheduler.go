package utils

import (
	"log"

	"github.com/robfig/cron"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
)

func startBatch() *cron.Cron {
	terminator := cron.New()

	terminator.AddFunc("0 * * * *", func() {
		if err := queries.IsExpiredFiles(); err != nil {
			log.Println("cron db query error", err.Error())
		}

		if err := queries.IsExpiredTexts(); err != nil {
			log.Println("cron db query error", err.Error())
		}
	})

	terminator.AddFunc("30 * * * *", func() {
		if err := queries.DelExpireFiles(); err != nil {
			log.Println("cron db query error", err.Error())
		}

		if err := queries.DelExpireTexts(); err != nil {
			log.Println("cron db query error", err.Error())
		}
	})
	terminator.Start()

	log.Println("Cron Job Start.... ✨")

	return terminator
}
