package main

import (
	http2 "ISEC/internal/api/delivery/http"
	repo "ISEC/internal/api/repository"
	"ISEC/internal/api/usecase"
	"ISEC/internal/middleware/proxy"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const proxyAddress = "0.0.0.0:8080"

func initDB() *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", "5432",
		"postgres", "password", "proxyDB")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {

	db := initDB()
	proxyRepo := repo.NewProxyRepo(db)
	proxy := &proxy.Proxy{ProxyRepo: proxyRepo}
	proxyUsecase := usecase.NewProxyUsecase(*proxyRepo, proxy)
	proxyHandler := http2.NewProxyHandler(*proxyUsecase)

	log.Println("Starting proxy server on", proxyAddress)
	go func() {
		if err := http.ListenAndServe(proxyAddress, proxy); err != nil {
			log.Fatal("Error during serving:", err)
		}
	}()

	mux := mux.NewRouter()
	mux.HandleFunc("/requests", proxyHandler.GetAllRequests)
	mux.HandleFunc("/request/{id}", proxyHandler.GetRequest)
	mux.HandleFunc("/repeat/{id}", proxyHandler.RepeatRequest)
	mux.HandleFunc("/scan/{id}", proxyHandler.ScanRequest)

	log.Println("Starting api server on", ":8000")
	log.Fatal(http.ListenAndServe(":8000", mux))

}
