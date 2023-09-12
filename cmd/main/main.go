package main

import (
	"ISEC/internal/middleware/proxy"
	"log"
	"net/http"
)

const proxyAddress = "127.0.0.1:8080"

func main() {
	handler := &proxy.Proxy{}

	log.Println("Starting proxy server on", proxyAddress)
	if err := http.ListenAndServe(proxyAddress, handler); err != nil {
		log.Fatal("Error during serving:", err)
	}
}
