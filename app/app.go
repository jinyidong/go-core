package app

import "go-core/app/server"

//查询实现了ServeHTTP接口的类
func init() {

}

func UseWebServer(config server.Config, handler server.Handler) server.Server {

	s := &server.WebServer{}

	s.Config = config

	s.Handler = handler

	return s
}
