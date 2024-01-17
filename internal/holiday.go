package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func todayIsHoliday() map[string]interface{} {
	// 获取当前日期
	now := time.Now().Format("2006-01-02")

	// 查询节假日是否存在
	query := "SELECT id, name, date, is_off_day FROM holiday WHERE date = ?"
	row := Db.QueryRow(query, now)
	var id int
	var name string
	var date string
	var isOffDay bool
	err := row.Scan(&id, &name, &date, &isOffDay)
	if err != nil {
		// 如果查询失败，则非节假日
		id = 0
	}

	if id == 0 {
		// 名称
		name := "工作日"
		// 是否为周末
		isOffDay := false
		// 判断是否为周末
		weekDay := time.Now().Weekday()
		if weekDay == time.Saturday || weekDay == time.Sunday {
			name = "周末"
			isOffDay = true
		}

		// 如果节假日不存在，则返回非节假日
		return map[string]interface{}{
			"name":      name,
			"date":      now,
			"isOffDay":  isOffDay,
			"isHoliday": isOffDay,
		}
	}

	// 返回节假日数据
	return map[string]interface{}{
		"name":      name,
		"date":      date,
		"isOffDay":  isOffDay,
		"isHoliday": true,
	}
}

func syncHolidayJob() {
	fmt.Println("开始同步节假日数据")

	// 获取当前年份
	year := time.Now().Year()
	// 今年
	nowYear := strconv.Itoa(year)
	// 明年
	nextYear := strconv.Itoa(year + 1)

	// 同步节假日数据
	syncHoliday(nowYear)  // 今年
	syncHoliday(nextYear) // 明年

	fmt.Println("同步节假日数据完成")
}

func syncHoliday(year string) {
	// 请求地址前缀
	prefixList := []string{
		"https://fastly.jsdelivr.net/gh/NateScarlet/holiday-cn@master/",
		"https://cdn.jsdelivr.net/gh/NateScarlet/holiday-cn@master/",
		"https://mirror.ghproxy.com/https://raw.githubusercontent.com/NateScarlet/holiday-cn/master/",
		"https://raw.githubusercontent.com/NateScarlet/holiday-cn/master/",
	}
	// 请求地址后缀
	suffix := ".json"

	var data []byte = nil
	// 请求数据 如果请求失败，则继续请求下一个地址。如果所有地址都请求失败，则不继续执行
	for _, prefix := range prefixList {
		url := prefix + year + suffix
		resp, err := http.Get(url)
		if err != nil {
			// 如果请求失败，则继续请求下一个地址
			continue
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			// 如果读取失败，则继续请求下一个地址
			continue
		}

		// 如果请求成功，则跳出循环
		break
	}

	// 如果请求失败，则不继续执行
	if data == nil {
		insertLog("HTTP 请求失败, Data is nil")
		fmt.Println("HTTP 请求失败, Data is nil")
		return
	}

	// 解析数据
	dataJson := make(map[string]interface{})
	err := json.Unmarshal(data, &dataJson)
	if err != nil {
		// 如果解析失败，则不继续执行
		insertLog("JSON 解析失败, Error: " + err.Error())
		fmt.Println("JSON 解析失败, Error: " + err.Error())
		return
	}

	// 获取节假日数据
	days := dataJson["days"].([]interface{})
	if days == nil {
		// 如果解析失败，则不继续执行
		insertLog("JSON 解析失败, Days is nil")
		return
	}

	// 如果节假日数据为空，则不继续执行
	if len(days) == 0 {
		insertLog("节假日数据为空, 跳过同步")
		return
	}

	// 遍历节假日数据
	for _, day := range days {
		// 节日名称
		name := day.(map[string]interface{})["name"].(string)
		// 节日日期
		date := day.(map[string]interface{})["date"].(string)
		// 是否为休息日
		isOffDay := day.(map[string]interface{})["isOffDay"].(bool)

		// 插入或更新节假日
		insertOrUpdateHoliday(name, date, isOffDay)
		fmt.Println("同步节假日成功, Name: " + name + ", Date: " + date + ", IsOffDay: " + strconv.FormatBool(isOffDay))
		insertLog("同步节假日成功, Name: " + name + ", Date: " + date + ", IsOffDay: " + strconv.FormatBool(isOffDay))
	}
}
