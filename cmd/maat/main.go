package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/nihei9/maat/service"
)

func main() {
	os.Exit(doMain())
}

func doMain() int {
	defer func() {
		log.Printf("shutdown")
	}()

	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	server := http.Server{
		Addr:    *addr,
		Handler: service.MakeHTTPHandler(),
	}

	log.Printf("listen on %v", *addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("%v", err)
	}

	return 0
}
