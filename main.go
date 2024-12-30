package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var isDev = flag.Bool("isDev", false, "Start the server in development environment")

func loadEnv() error {
	if err := godotenv.Load(".env.local"); err != nil {
		return fmt.Errorf("failed to load env file %v", err)
	}
	return nil
}

func main() {
	flag.Parse()
	prefix := "[main]"

	if bool(*isDev) {
		log.Println(prefix+ " loading environmental variables from .env.local")
		err := loadEnv()

		if err != nil {
			log.Fatalf("%s %v",prefix, err)
		}
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	listenAddr := host + ":" + port

	if host == "" || port == "" {
		log.Fatalf("%s Environmental variables not provided", prefix)
	}
	router := mux.NewRouter()

	RegisterRoutes(router)
	
	started := make(chan struct{})
	go func() {
		log.Println(prefix, "Starting server on ")
		if err := http.ListenAndServe(listenAddr, router); err != nil {
			log.Fatalf(prefix+" Server failed: %v", err)
		}
	}()
	
	go func() {
		for {
			resp, err := http.Get("/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				close(started)
				return
			}
			started <- struct{}{}
		}
	}()

	<-started
	log.Println(prefix, "Application startup complete")
	
	select {}
}