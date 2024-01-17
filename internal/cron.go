package internal

import "github.com/robfig/cron/v3"

func InitCron() {
	// 执行前先同步一次节假日数据
	syncHoliday()

	// 创建定时任务
	c := cron.New()
	// 每月 25 日 01:00 执行
	_, err := c.AddFunc("0 0 1 25 * *", syncHoliday)
	if err != nil {
		// 如果添加失败，则不继续执行
		insertLog("添加定时任务失败, Error: " + err.Error())
		return
	}
	// 启动定时任务
	c.Start()
}
