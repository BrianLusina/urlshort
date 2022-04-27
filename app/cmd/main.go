package main

import (
	"flag"
	"net/http"
	"net/rpc"
	"urlshort/app/api"
	"urlshort/app/internal/core/contracts"
	"urlshort/app/internal/repo"
)

var (
	primaryAddr = flag.String("primary", "", "RPC address of the primary store")
	listenAddr  = flag.String("http", ":3000", "http listen address")
	dataFile    = flag.String("file", "store.json", "data store file name")
	hostname    = flag.String("host", "localhost", "host name and port")
	rpcEnabled  = flag.Bool("rpc", false, "enable rpc server")
)

var store contracts.Repository

func main() {
	flag.Parse()

	if *primaryAddr != "" {
		store = repo.NewProxyStore(*primaryAddr)
	} else {
		store = repo.NewUrlStore(*dataFile)
	}

	handler := api.NewHandler(store)

	if *rpcEnabled {
		rpc.RegisterName("Store", store)
		rpc.HandleHTTP()
	}

	http.HandleFunc("/", handler.Redirect)
	http.HandleFunc("/add", handler.Add)
	http.ListenAndServe(*listenAddr, nil)
}
