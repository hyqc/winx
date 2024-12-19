package wnet

import (
	"fmt"
	"time"
)

type LogLevelName string

const (
	LogDebugName LogLevelName = "DEBUG"
	LogInfoName  LogLevelName = "INFO"
	LogWarnName  LogLevelName = "WARN"
	LogErrorName LogLevelName = "ERROR"
	LogPanicName LogLevelName = "PANIC"
)

func Print(args ...any) {
	fmt.Println("============================: ", args)
}

func SysPrintDebug(args ...any) {
	fmt.Println(SysPrintLogLevel(LogDebugName), args)
}

func SysPrintInfo(args ...any) {
	fmt.Println(SysPrintLogLevel(LogInfoName), args)
}

func SysPrintWarn(args ...any) {
	fmt.Println(SysPrintLogLevel(LogWarnName), args)
}

func SysPrintError(args ...any) {
	fmt.Println(SysPrintLogLevel(LogErrorName), args)
}

func SysPrintPanic(args ...any) {
	panic(fmt.Errorf(SysPrintLogLevel(LogPanicName), args))
}

func SysPrintLogLevel(name LogLevelName) string {
	return fmt.Sprintf("[SERVER] [%s] [%s] [MSG]: ", name, time.Now().Format(time.RFC3339))
}
