package main

import (
	"log"
	"net/http"
	"os"

	h "github.com/mehrdadrad/servicemesh-examples/pkg/http"
)

type service struct {
	addr    string
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
	addr := os.Getenv("FRONTEND_ADDR")
	if addr == "" {
		addr = ":8082"
	}

	backend := os.Getenv("MIDDLEWARE")
	if backend == "" {
		log.Fatal("middleware address is not available")
	}

	return &service{
		addr:    addr,
		backend: backend,
	}
}

func (s service) start() {
	log.Printf("frontend is starting %s", s.addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/time", getTime(s.backend))
	log.Println(http.ListenAndServe(s.addr, mux))
}

func getTime(backend string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := h.HTTPClient("get", backend)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(code)
			w.Write([]byte("current time "))
			w.Write(resp)
		}
	}
}
