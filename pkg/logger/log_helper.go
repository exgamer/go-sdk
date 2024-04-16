package logger

import (
	"github.com/exgamer/go-sdk/pkg/config"
	"github.com/exgamer/go-sdk/pkg/exception"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Info Обычный лог
func Info(format string, v ...any) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	infoLog.Printf(format, v)
}

// Error лог с ошибкой
func Error(format string, v ...any) {
	infoLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	infoLog.Printf(format, v)
}

// LogError лог с ошибкой
func LogError(err error) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog.Println(err)
}

// LogAppException лог AppException
func LogAppException(appException *exception.AppException) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog.Println(appException.Error.Error())
}

// FormattedInfo Форматированный лог
func FormattedInfo(serviceName string, method string, uri string, status int, requestId string, message string) {
	FormattedLog("INFO", serviceName, method, uri, status, requestId, message)
}

// FormattedError Форматированный лог ошибки
func FormattedError(serviceName string, method string, uri string, status int, requestId string, message string) {
	FormattedLog("ERROR", serviceName, method, uri, status, requestId, message)
}

// FormattedLog Форматированный лог
func FormattedLog(level string, serviceName string, method string, uri string, status int, requestId string, message string) {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	messageBuilder := strings.Builder{}
	messageBuilder.WriteString(time.Now().Format("2006-01-02 15:04:05.345"))
	messageBuilder.WriteString(" " + level + " ")
	messageBuilder.WriteString("[" + serviceName + "," + requestId + "]")
	messageBuilder.WriteString("[" + method + "," + uri + "," + strconv.Itoa(status) + "]")
	messageBuilder.WriteString(" " + message)

	log.Println(messageBuilder.String())
	log.SetFlags(log.Ldate | log.Ltime)
}

// FormattedLogWithAppInfo Форматированный лог для RequestData
func FormattedLogWithAppInfo(appInfo *config.AppInfo, message string) {
	FormattedInfo(appInfo.ServiceName, appInfo.RequestMethod, appInfo.RequestUrl, 0, appInfo.RequestId, message)
}

// FormattedErrorWithAppInfo Форматированный лог ошибки для RequestData
func FormattedErrorWithAppInfo(appInfo *config.AppInfo, message string) {
	FormattedInfo(appInfo.ServiceName, appInfo.RequestMethod, appInfo.RequestUrl, 1, appInfo.RequestId, message)
}
