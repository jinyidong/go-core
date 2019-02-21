package log

import (
	"testing"
	"time"
)

func TestLog_Save(t *testing.T) {
	var log = &Log{
		Message:    "hello",
		Exception:  "exception",
		Level:      "DEBUG",
		LogType:    "MonitorApi",
		Logger:     "MonitorApi",
		Userip:"::",
		Timespan:time.Now().Format(time.RFC3339),
		CreateTime: time.Now().Format("2006-01-02 15:04:05") + "+08:00",
	}

	log.Save()
}

func BenchmarkLog_Save(b *testing.B) {
	for i:=0;i<b.N;i++ {
		var log = &Log{
			Message:    "hello",
			Exception:  "exception",
			Level:      "DEBUG",
			LogType:    "MonitorApi",
			Logger:     "MonitorApi",
			Userip:     "::",
			Timespan:   time.Now().Format(time.RFC3339),
			CreateTime: time.Now().Format("2006-01-02 15:04:05") + "+08:00",
		}
		log.Save()
	}
}

func TestDefaultLogger_Info(t *testing.T) {
	DefaultLogger.Info("MonitorApi","MonitorApi","MonitorApiMessage",1000,"123","::","")

	time.Sleep(5*time.Second)
}

func TestDefaultLogger_Error(t *testing.T) {
	DefaultLogger.Error("MonitorApi","MonitorApi","message","123","::","","")
	time.Sleep(10*time.Second)
}

func TestDefaultLogger_Warn(t *testing.T) {
	DefaultLogger.Warn("MonitorApi","MonitorApi","message","404","::","","")
	time.Sleep(10*time.Second)
}