package app

import (
	"fmt"
	"sygit.suiyi.com.cn/go/core/app/rpc"
	"sygit.suiyi.com.cn/go/core/app/server"
	"sygit.suiyi.com.cn/go/core/http"
	"sygit.suiyi.com.cn/go/core/log"
	"sygit.suiyi.com.cn/go/core/mqbus"
	"sygit.suiyi.com.cn/go/core/util"
	"testing"
	"time"
)

type logHandler struct {
}

func (*logHandler) Handle(url string, data []byte) ([]byte, error) {
	var err error
	var r []byte

	if url == "/api/logs/eventlog/create" {
		l := &log.Log{}
		err = util.ByteToStruct(data, l)
		fmt.Println(util.StructToJson(l))
		if err != nil {
			return nil, err
		}
		go func(log *log.Log) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(time.Now().Format("2006-01-02 15:04:05 +08:00"), "LogHandler", r)
				}
			}()
			mqbus.Default.Post(log)
		}(l)

		if err != nil {
			return nil, err
		} else {
			r = []byte("ok")
		}
	}

	return r, nil
}

func TestUseRpcServer(t *testing.T) {
	UseRpcServer(server.Config{
		Port:    8080,
		TimeOut: 30,
	}, &logHandler{}).Run()
}

func TestLog(t *testing.T) {
	client := rpc.NewClient("127.0.0.1:8080")

	defer client.Close()

	finish := make(chan struct{})
	json := fmt.Sprintf(`{"traceId":"9adb0bbc-b809-46af-9737-9e7e1d62cb8e","logType":"MonitorApi","logger":"MonitorApi11111111111","level":"INFO","createTime":"%s","msec":0,"userip":"10.1.4.137","message":"2018-09-28 17:08:53.090 [main] INFO  c.s.a.l.c.t.LogTest - test info","exception":"null"}`, time.Now().Format("2006-01-02 15:04:05 +08:00"))
	client.Call("/api/logs/eventlog/create", json, func(res []byte) {
		finish <- struct{}{}
	})
	<-finish
}

func BenchmarkUseRpcServer(b *testing.B) {

	client := rpc.NewClient("127.0.0.1:8080")

	defer client.Close()

	for i := 0; i < b.N; i++ {
		finish := make(chan struct{})
		json := fmt.Sprintf(`{"traceId":"9adb0bbc-b809-46af-9737-9e7e1d62cb8e","logType":"MonitorApi","logger":"MonitorApi11111111111","level":"INFO","createTime":"%s","msec":0,"userip":"10.1.4.137","message":"2018-09-28 17:08:53.090 [main] INFO  c.s.a.l.c.t.LogTest - test info","exception":"null"}`, time.Now().Format("2006-01-02 15:04:05 +08:00"))
		client.Call("/api/logs/eventlog/create", json, func(res []byte) {
			finish <- struct{}{}
		})
		<-finish
	}
}

func TestUseWebServer(t *testing.T) {
	UseWebServer(server.Config{
		Port:    8080,
		TimeOut: 30,
	}, &logHandler{}).Run()
}

func TestUseWebServer_Log(t *testing.T) {

	json := fmt.Sprintf(`{"traceId":"9adb0bbc-b809-46af-9737-9e7e1d62cb8e","logType":"MonitorApi","logger":"MonitorApi11111111111","level":"INFO","createTime":"%s","msec":0,"userip":"10.1.4.137","message":"2018-09-28 17:08:53.090 [main] INFO  c.s.a.l.c.t.LogTest - test info","exception":"null"}`, time.Now().Format("2006-01-02 15:04:05 +08:00"))

	response, err := http.PostJson("http://127.0.0.1:8080/api/logs/eventlog/create", json)

	if err == nil {
		fmt.Println(string(response))
	}
}
