package app

import "github.com/jinyidong/go-core/app/server"

func UseWebServer(config server.Config, handler server.Handler) server.Server {
	s := &server.WebServer{}

	s.Config = config

	s.Handler = handler

	return s
}
