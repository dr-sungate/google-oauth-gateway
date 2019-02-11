package logger

import (
	"fmt"
	"github.com/dr-sungate/google-oauth-gateway/api/service/utils"
	"log"
	"runtime"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

func Debug(msg ...interface{}) {
	if utils.GetEnv("VERIFY_MODE", "") == "enable" {
		_, filename, line, _ := runtime.Caller(1)
		log.Println(fmt.Sprintf("[%s] %s:%d ", "DEBG", filename, line), msg)
	}
}

func Info(msg ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	log.Println(fmt.Sprintf("[%s] %s:%d ", "INFO", filename, line), msg)
}

func Warn(msg ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	log.Println(fmt.Sprintf("[%s] %s:%d ", "WARN", filename, line), msg)
}

func Error(msg ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	log.Println(fmt.Sprintf("[%s] %s:%d ", "ERROR", filename, line), msg)
}
