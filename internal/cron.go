package internal

import "github.com/robfig/cron/v3"

func InitCron() {
	// 每月 25 日 01:00 执行
	c := cron.New()
	_, err := c.AddFunc("0 0 1 25 * *", syncHoliday)
	if err != nil {
		// 如果添加失败，则不继续执行
		insertLog("添加定时任务失败, Error: " + err.Error())
		return
	}
	c.Start()
}
