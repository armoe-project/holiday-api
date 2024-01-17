package main

import "zhenxin.me/holiday/internal"

func main() {
	internal.InitDatabase()
	internal.InitCron()
	internal.InitHttpServer()
}
