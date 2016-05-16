package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
)

func main() {
	srv := rpc.NewServer()
	srv.RegisterCodec(json.NewCodec(), "application/json")
	srv.RegisterService(Service{}, "Service")
	http.Handle("/rpc", srv)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("listen: %s", err)
	}
}

type Service struct{}

func (_ Service) Sum(_ *http.Request, req *[]int, sum *int) error {
	for _, n := range *req {
		*sum += n
	}
	return nil
}

func (_ Service) Echo(_ *http.Request, req, res *interface{}) error {
	*res = *req
	return nil
}
