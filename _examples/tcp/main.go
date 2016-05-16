package main

import (
	"encoding/json"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	srv := rpc.NewServer()
	srv.Register(Service{})
	listen(srv)
}

func listen(srv *rpc.Server) {
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("listen: %s", err)
	}

	for {
		sock, err := l.Accept()
		if err != nil {
			log.Fatalf("accept: %s", err)
		}

		go srv.ServeCodec(jsonrpc.NewServerCodec(sock))
	}
}

type Service struct{}

func (_ Service) Sum(req *[]int, sum *int) error {
	for _, n := range *req {
		*sum += n
	}
	return nil
}

func (_ Service) Echo(req, res *json.RawMessage) error {
	*res = *req
	return nil
}
