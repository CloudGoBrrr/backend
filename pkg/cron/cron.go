package cron

import (
	"cloudgobrrr/backend/database/model"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func InitCron() error {
	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		fmt.Println("[Daily Cleanup] " + time.Now().Format("02-01-2006"))
		model.DownloadSecretCleanup()
		model.SessionCleanup()
	})
	c.Start()
	return nil
}
