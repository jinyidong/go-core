package server

import "net/http"

type WebServer struct {
}

type Handler func(http.ResponseWriter, *http.Request)

func Get(path string, handler Handler) {

}
