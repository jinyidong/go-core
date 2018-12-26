package log

import (
	"fmt"
	"github.com/jinyidong/go-core/mqbus"
	"sync"
	"time"
)

type Log struct {
	Message        interface{} `json:"message"`
	Msec           int64       `json:"msec"`
	Userip         string      `json:"userip"`
	Exception      interface{} `json:"exception"`
	Timespan       string      `json:"timespan"`
	Level          string      `json:"level"`
	Logger         string      `json:"logger"`
	CreateTime     string      `json:"createTime"`
	LogType        string      `json:"logType"`
	TraceId        string      `json:"traceId"`
	CollectionName string      `json:"collectionName"`
}

func (*Log) ConnectionString() string {
	return "logmq/1"
}
func (*Log) Exchange() string {
	return "Core.Logs.EventLog:Core"
}
func (*Log) MessageType() string {
	return "Core.Logs.EventLog:Core"
}
func (*Log) RoutingKey() string {
	return "#"
}

func (l *Log) Save() {

	go func(log *Log) {

		err := mqbus.Default.Post(log)
		if err != nil {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05 +08:00"), "LOG::ERROR::", err)
		}

	}(l)
}

type Logger interface {
	Info(logType string, logger string, message string, msec int64, traceId string, clientIp string, collectionName string)
	Error(logType string, logger string, exception string, message string, traceId string, clientIp string, collectionName string)
	Warn(logType string, logger string, exception string, message string, traceId string, clientIp string, collectionName string)
}

var DefaultLogger Logger = newDefaultLogger()

type defaultLogger struct {
	entries []*Log
	mu      sync.Mutex
}

func newDefaultLogger() *defaultLogger {
	logger := &defaultLogger{}
	go func() {
		for {
			if len(logger.entries) > 0 {
				logger.mu.Lock()
				for _, entry := range logger.entries {
					entry.Save()
				}
				logger.entries = nil
				logger.mu.Unlock()
			} else {
				time.Sleep(5 * time.Second)
			}
		}

	}()
	return logger
}

func (l *defaultLogger) Info(logType string, logger string, message string, msec int64, traceId string, clientIp string, collectionName string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log := &Log{
		Message:        message,
		Level:          "INFO",
		LogType:        logType,
		Logger:         logger,
		Userip:         clientIp,
		Msec:           msec,
		Timespan:       time.Now().Format(time.RFC3339),
		CreateTime:     time.Now().Format(time.RFC3339),
		TraceId:        traceId,
		CollectionName: collectionName,
	}

	l.entries = append(l.entries, log)
}

func (l *defaultLogger) Error(logType string, logger string, exception string, message string, traceId string, clientIp string, collectionName string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log := &Log{
		Exception:      exception,
		Level:          "ERROR",
		Message:        message,
		LogType:        logType,
		Logger:         logger,
		Userip:         clientIp,
		TraceId:        traceId,
		Timespan:       time.Now().Format(time.RFC3339),
		CreateTime:     time.Now().Format("2006-01-02 15:04:05") + "+08:00",
		CollectionName: collectionName,
	}

	l.entries = append(l.entries, log)
}

func (l *defaultLogger) Warn(logType string, logger string, exception string, message string, traceId string, clientIp string, collectionName string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log := &Log{
		Exception:      exception,
		Level:          "WARN",
		Message:        message,
		LogType:        logType,
		Logger:         logger,
		Userip:         clientIp,
		TraceId:        traceId,
		Timespan:       time.Now().Format(time.RFC3339),
		CreateTime:     time.Now().Format("2006-01-02 15:04:05") + "+08:00",
		CollectionName: collectionName,
	}

	l.entries = append(l.entries, log)
}
