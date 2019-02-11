package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

func Debug(msg ...interface{}) {
	if os.Getenv("VERIFY_MODE") == "enable" {
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
