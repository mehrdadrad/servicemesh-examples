package main

import (
	"log"
	"net"
	"net/http"
	"os"

	h "github.com/mehrdadrad/servicemesh-examples/pkg/http"
)

type service struct {
	port string
}

func main() {
	s := newService()
	s.start()
}

func newService() *service {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	return &service{
		port: port,
	}
}

func (s service) start() {
	log.Println("backend is starting")
	mux := http.NewServeMux()
	mux.HandleFunc("/time", getTime)
	log.Println(http.ListenAndServe(net.JoinHostPort("", s.port), mux))
}

func getTime(w http.ResponseWriter, r *http.Request) {
	resp, err := h.HttpClient("get", "http://time.jsontest.com")
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(resp)
	}
}
