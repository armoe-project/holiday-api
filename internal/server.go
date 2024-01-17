package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func InitHttpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 结果 JSON
		result := make(map[string]interface{})
		result["code"] = 0
		result["message"] = "success"
		result["data"] = todayIsHoliday()

		// 返回结果
		resultJson, err := json.Marshal(result)
		if err != nil {
			insertLog("JSON 转换失败, Error: " + err.Error())
			fmt.Println("JSON 转换失败, Error: " + err.Error())
			return
		}
		_, err = w.Write(resultJson)
		if err != nil {
			insertLog("HTTP 写入失败, Error: " + err.Error())
			fmt.Println("HTTP 写入失败, Error: " + err.Error())
			return
		}
	})

	// 从环境变量中获取端口号
	port := os.Getenv("PORT")
	if port == "" {
		// 如果环境变量中没有端口号，则使用默认端口号
		port = "9527"
	}
	fmt.Println("节假日同步服务已在端口 " + port + " 启动")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		insertLog("HTTP 服务启动失败, Error: " + err.Error())
		fmt.Println("HTTP 服务启动失败, Error: " + err.Error())
		return
	}
}
