package main

import (
	"log"
	"net/http"
	"os"

	h "github.com/mehrdadrad/servicemesh-examples/pkg/http"
)

type service struct {
	addr string
}

func main() {
	s := newService()
	s.start()
}

func newService() *service {
	addr := os.Getenv("BACKEND_ADDR")
	if len(addr) < 2 {
		addr = ":8080"
	}

	return &service{
		addr: addr,
	}
}

func (s service) start() {
	log.Printf("backend is starting %s", s.addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/time", getTime)
	log.Println(http.ListenAndServe(s.addr, mux))
}

func getTime(w http.ResponseWriter, r *http.Request) {
	resp, err := h.HttpClient("get", "http://time.jsontest.com")
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(resp)
	}
}
