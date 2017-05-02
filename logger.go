package tcpsocket

import (
	"fmt"
	"runtime"
	"time"
)

func traceFile() (file string, line int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	return file, line
}

func outputLog(levelname string, log string) string {
	file, line := traceFile()
	vCurrTime := time.Now()
	return fmt.Sprintf("%.4d-%.2d-%.2d %.2d:%.2d:%.2d.%.6d[%s:%d][%s]%s",
		vCurrTime.Year(), vCurrTime.Month(), vCurrTime.Day(),
		vCurrTime.Hour(), vCurrTime.Minute(), vCurrTime.Second(),
		vCurrTime.Nanosecond()/1000, file, line, levelname, log)
}

func Debugln(v ...interface{}) {
	fmt.Print(outputLog(`debug`, fmt.Sprintln(v...)))
}

func Debugf(f string, v ...interface{}) {
	fmt.Print(outputLog(`debug`, fmt.Sprintf(f, v...)))
}

func Errorln(v ...interface{}) {
	fmt.Print(outputLog(`error`, fmt.Sprintln(v...)))
}

func Errorf(f string, v ...interface{}) {
	fmt.Print(outputLog(`error`, fmt.Sprintf(f, v...)))
}
