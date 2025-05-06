package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sithu-go/go-soa/common/models"
)

var productApiURL, userApiURL string

var orders = []models.Order{}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, order := range orders {
		if order.ID == params["id"] {
			json.NewEncoder(w).Encode(order)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// Validate user exists
	resp, err := http.Get(userApiURL + "/users/" + order.UserID + "/validate")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to validate user"})
		return
	}
	defer resp.Body.Close()

	var userValid map[string]bool
	json.NewDecoder(resp.Body).Decode(&userValid)

	if !userValid["valid"] {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User does not exist"})
		return
	}

	// Process each product in the order
	var total float64 = 0
	for _, item := range order.Products {
		// Get product details
		resp, err := http.Get(productApiURL + "/products/" + item.ProductID)
		if err != nil || resp.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Product not found: " + item.ProductID})
			return
		}

		var product models.Product
		json.NewDecoder(resp.Body).Decode(&product)
		resp.Body.Close()

		// Update the price in the order item from the product service
		// order.Products[i].Price = product.Price

		// Add to total
		total += product.Price * float64(item.Quantity)

		// Update stock
		stockUpdate := map[string]int{"quantity": item.Quantity}
		jsonData, _ := json.Marshal(stockUpdate)

		stockReq, _ := http.NewRequest("PUT", productApiURL+"/products/"+item.ProductID+"/stock", bytes.NewBuffer(jsonData))
		stockReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		stockResp, err := client.Do(stockReq)

		if err != nil || stockResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(stockResp.Body)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update stock: " + string(body)})
			stockResp.Body.Close()
			return
		}
		stockResp.Body.Close()
	}

	// Set order details
	order.ID = uuid.New().String()
	order.Total = total
	order.Status = "created"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	// Save order
	orders = append(orders, order)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func main() {
	env := os.Getenv("RUNNING_ENV")
	isDocker := env == "docker"

	if isDocker {
		productApiURL = "http://product-service:8081"
		userApiURL = "http://user-service:8082"
	} else {
		productApiURL = "http://localhost:8081"
		userApiURL = "http://localhost:8082"
	}

	r := mux.NewRouter()

	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")

	log.Println("Order Service is running on port 8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}
