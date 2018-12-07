package server

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
