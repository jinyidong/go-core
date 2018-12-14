package mqbus

import (
	"testing"
	"time"
	"fmt"
)

type log struct {
	Message    interface{} `json:"message"`
	Msec       int `json:"msec"`
	Userip     string `json:"userip"`
	Exception  interface{} `json:"exception"`
	Timespan   string `json:"timespan"`
	Level      string `json:"level"`
	Logger     string `json:"logger"`
	CreateTime string `json:"createTime"`
	LogType    string `json:"logType"`
	TraceId    string `json:"traceId"`
}

func (*log) ConnectionString() string {
	////connectionString := "host=10.1.4.131:5672;username=guest;password=guest"
	//connectionString := "host=10.1.62.23:5672;username=shampoo;password=123456"
	return "logmq/1"
}
func (*log) Exchange() string {
	return "Core.Logs.EventLog:Core"
}
func (*log) MessageType() string {
	return "Core.Logs.EventLog:Core"
}
func (*log) RoutingKey() string {
	return "#"
}

func TestMqBus_Post(t *testing.T) {

	var log = &log{
		Message:    "hello",
		Exception:  "exception",
		Level:      "DEBUG",
		LogType:    "MonitorApi",
		Logger:     "MonitorApi",
		CreateTime: time.Now().Format("2006-01-02 15:04:05") + "+08:00",
	}

	err := Default.Post(log)

	if err != nil {
		t.Fatalf(err.Error())
	}

}


func TestMqBus_DelayPost(t *testing.T) {

	for {
		var log = &log{
			Message:    "hello",
			Exception:  "exception",
			Level:      "DEBUG",
			LogType:    "MonitorApi",
			Logger:     "MonitorApi",
			CreateTime: time.Now().Format("2006-01-02 15:04:05") + "+08:00",
		}

		err := Default.Post(log)

		if err != nil {
			fmt.Println(err)
		}else {
			fmt.Println("log")
		}

		time.Sleep(time.Second *20)
	}
}

func BenchmarkMqBus_Post(b *testing.B) {
	for i:=0;i<b.N;i++{
		var log = &log{
			Message:    "hello",
			Exception:  "exception",
			Level:      "DEBUG",
			LogType:    "MonitorApi",
			Logger:     "MonitorApi",
			CreateTime: time.Now().Format("2006-01-02 15:04:05") + "+08:00",
		}

		err := Default.Post(log)

		if err != nil {
			b.Fatalf(err.Error())
		}
	}
}
