package internal

import "github.com/robfig/cron/v3"

func InitCron() {
	// 执行前先同步一次节假日数据
	syncHolidayJob()

	// 创建定时任务
	c := cron.New()
	// 每日 00:00:00 执行一次
	_, err := c.AddFunc("@daily", syncHolidayJob)
	if err != nil {
		// 如果添加失败，则不继续执行
		insertLog("添加定时任务失败, Error: " + err.Error())
		return
	}

	// 启动定时任务
	c.Start()
}
