package main

import (
	"log"
	"net"
	"net/http"
	"os"

	h "github.com/mehrdadrad/servicemesh-examples/pkg/http"
)

type service struct {
	port    string
	backend string
}

type time struct {
	Date  string `json:"date"`
	Time  string `json:"time"`
	Epoch int64  `json:"milliseconds_since_epoch"`
}

func main() {
	s := newService()
	s.start()
}

func newService() *service {
	port := os.Getenv("FRONTEND_PORT")
	if port == "" {
		port = "8082"
	}

	backend := os.Getenv("MIDDLEWARE")
	if backend == "" {
		log.Fatal("middleware address is not available")
	}

	return &service{
		port:    port,
		backend: backend,
	}
}

func (s service) start() {
	log.Println("frontend is starting")
	mux := http.NewServeMux()
	mux.HandleFunc("/time", getTime(s.backend))
	log.Println(http.ListenAndServe(net.JoinHostPort("", s.port), mux))
}

func getTime(backend string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := h.HttpClient("get", backend)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("current time "))
			w.Write(resp)
		}
	}
}
