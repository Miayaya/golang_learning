package middleware

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	// 确保 logs 文件夹存在
	logDir := "./logs"
	_ = os.MkdirAll(logDir, 0755)

	// 按日期生成日志文件
	logName := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("create log file err:", err)
	} else {
		// 同时输出到控制台和文件
		log.SetOutput(os.Stdout)
		gin.DefaultWriter = file
		gin.DefaultErrorWriter = file
	}

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return param.TimeStamp.Format("2006-01-02 15:04:05") +
			" | " + param.Method +
			" | " + param.Path +
			" | " + param.ClientIP +
			" | " + strconv.Itoa(param.StatusCode) +
			" | " + param.Latency.String() + "\n"
	})
}
