package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
)

type Config struct {
	Port    int64
	TimeOut int64
}

type Server interface {
	Run()
}

type Handler interface {
	Handle(uri string, data []byte) ([]byte, error)
}

type WebServer struct {
	Config  Config
	Handler Handler
	//ServeMux map[string]http.Handler
}

func (s *WebServer) Run() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("listening .....:", s.Config.Port)

	http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port), &echo{server: s})
}

type echo struct {
	server *WebServer
}

func (h echo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := path.Clean(req.URL.Path)

	if req.Body == nil {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	buf, err := h.server.Handler.Handle(url, body)

	w.Write(buf)
}
