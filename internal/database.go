package internal

import (
	"database/sql"
	"fmt"
)

import _ "github.com/mattn/go-sqlite3"

func InitDatabase() {
	// 连接数据库
	dsn := "holiday.db"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		// 如果连接失败，则不继续执行
		fmt.Println("数据库连接失败, Error: " + err.Error())
		return
	}
	Db = db
	initTable()
}

func initTable() {
	createHolidayTable()
	createLogTable()
}

func createHolidayTable() {
	// 创建节假日表
	create := "CREATE TABLE IF NOT EXISTS holiday (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR(255), date VARCHAR(255), is_off_day BOOLEAN)"
	_, err := Db.Exec(create)
	if err != nil {
		// 如果创建失败，则不继续执行
		fmt.Println("创建数据库表失败, Error: " + err.Error())
		return
	}
}

func createLogTable() {
	// 创建日志表
	create := "CREATE TABLE IF NOT EXISTS log (id INTEGER PRIMARY KEY AUTOINCREMENT, content VARCHAR(255), create_time DATETIME)"
	_, err := Db.Exec(create)
	if err != nil {
		// 如果创建失败，则不继续执行
		fmt.Println("创建数据库表失败, Error: " + err.Error())
		return
	}
}

func insertOrUpdateHoliday(name string, date string, isOffDay bool) {
	// 查询节假日是否存在
	query := "SELECT id FROM holiday WHERE date = ?"
	row := Db.QueryRow(query, date)
	var id int
	err := row.Scan(&id)
	if err != nil {
		// 如果查询失败，则不继续执行
		fmt.Println("查询节假日失败, Error: " + err.Error())
		return
	}

	if id == 0 {
		// 如果不存在，则插入
		insert := "INSERT INTO holiday (name, date, is_off_day) VALUES (?, ?, ?)"
		_, err := Db.Exec(insert, name, date, isOffDay)
		if err != nil {
			// 如果插入失败，则不继续执行
			fmt.Println("插入节假日失败, Error: " + err.Error())
			return
		}
	} else {
		// 如果存在，则更新
		update := "UPDATE holiday SET name = ?, is_off_day = ? WHERE id = ?"
		_, err := Db.Exec(update, name, isOffDay, id)
		if err != nil {
			// 如果更新失败，则不继续执行
			fmt.Println("更新节假日失败, Error: " + err.Error())
			return
		}
	}
}

func insertLog(content string) {
	// 插入日志
	insert := "INSERT INTO log (content, create_time) VALUES (?, datetime('now', 'localtime'))"
	_, err := Db.Exec(insert, content)
	if err != nil {
		// 如果插入失败，则不继续执行
		fmt.Println("插入日志失败, Error: " + err.Error())
		return
	}
}

var Db *sql.DB
