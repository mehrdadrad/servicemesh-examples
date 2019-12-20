package main

import (
	"encoding/json"
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
	port := os.Getenv("MIDDLEWARE_PORT")
	if port == "" {
		port = "8081"
	}

	backend := os.Getenv("BACKEND")
	if backend == "" {
		log.Fatal("backend address is not available")
	}

	return &service{
		port:    port,
		backend: backend,
	}
}

func (s service) start() {
	log.Println("middleware is starting")
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
			t := &time{}
			err := json.Unmarshal(resp, t)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			w.Write([]byte(t.Date + " " + t.Time))
		}
	}
}
