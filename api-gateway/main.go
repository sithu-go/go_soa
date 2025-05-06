package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

// Simple API Gateway to route requests to appropriate services
func main() {
	env := os.Getenv("RUNNING_ENV")
	isDocker := env == "docker"

	var productApiURL string
	var userApiURL string
	var orderApiURL string

	if isDocker {
		productApiURL = "http://product-service:8081"
		userApiURL = "http://user-service:8082"
		orderApiURL = "http://order-service:8083"
	} else {
		productApiURL = "http://localhost:8081"
		userApiURL = "http://localhost:8082"
		orderApiURL = "http://localhost:8083"
	}

	r := mux.NewRouter()

	// Route to Product Service
	productServiceURL, err := url.Parse(productApiURL)
	if err != nil {
		log.Fatal(err)
	}
	productProxy := httputil.NewSingleHostReverseProxy(productServiceURL)
	r.PathPrefix("/products").Handler(productProxy)

	// Route to User Service
	userServiceURL, err := url.Parse(userApiURL)
	if err != nil {
		log.Fatal(err)
	}
	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	r.PathPrefix("/users").Handler(userProxy)

	// Route to Order Service
	orderServiceURL, err := url.Parse(orderApiURL)
	if err != nil {
		log.Fatal(err)
	}
	orderProxy := httputil.NewSingleHostReverseProxy(orderServiceURL)
	r.PathPrefix("/orders").Handler(orderProxy)

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API Gateway is healthy"))
	}).Methods("GET")

	log.Println("API Gateway is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
