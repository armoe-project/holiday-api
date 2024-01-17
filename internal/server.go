package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func InitHttpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 输出日志
		method := r.Method         // 请求方法
		uri := r.RequestURI        // 请求 URI
		remoteAddr := r.RemoteAddr // 请求 IP
		// 如果请求头中有 X-Forwarded-For 或 X-Real-IP，则使用请求头中的 IP
		if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
			remoteAddr = ip
		} else if ip := r.Header.Get("X-Real-IP"); ip != "" {
			remoteAddr = ip
		}
		userAgent := r.UserAgent()                    // 请求 UA
		t := time.Now().Format("2006-01-02 15:04:05") // 当前时间

		fmt.Println("[" + t + "] " + remoteAddr + " " + method + " " + uri + " " + userAgent)

		// 结果 JSON
		result := make(map[string]interface{})
		result["code"] = 0
		result["message"] = "success"
		result["data"] = nil

		// 获取请求参数
		date := r.URL.Query().Get("date")

		// 如果参数不为空，则校验是否为日期格式
		if date != "" {
			_, err := time.Parse("2006-01-02", date)
			if err != nil {
				// 如果不是日期格式，则返回错误
				result["code"] = 1
				result["message"] = "参数错误"
			} else {
				// 查询节假日
				result["data"] = todayIsHoliday(date)
			}
		} else {
			// 查询节假日
			result["data"] = todayIsHoliday(date)
		}

		// 返回结果
		resultJson, err := json.Marshal(result)
		if err != nil {
			insertLog("JSON 转换失败, Error: " + err.Error())
			fmt.Println("JSON 转换失败, Error: " + err.Error())
			// 如果转换失败，自行拼接 JSON
			resultJson = []byte(`{"code": 1, "message": "JSON 转换失败", "data": null}`)
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
