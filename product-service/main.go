package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sithu-go/go-soa/common/models"
)

var products = []models.Product{}

func init() {
	// Add some sample products
	products = append(products, models.Product{
		ID:          "1",
		Name:        "Laptop",
		Description: "Powerful laptop with latest processor",
		Price:       999.99,
		StockCount:  10,
	})
	products = append(products, models.Product{
		ID:          "2",
		Name:        "Smartphone",
		Description: "Latest smartphone with great camera",
		Price:       499.99,
		StockCount:  15,
	})
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range products {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	product.ID = uuid.New().String()
	products = append(products, product)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func updateStockCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var stockUpdate struct {
		Quantity int `json:"quantity"`
	}

	err := json.NewDecoder(r.Body).Decode(&stockUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	for i, product := range products {
		if product.ID == params["id"] {
			if products[i].StockCount >= stockUpdate.Quantity {
				products[i].StockCount -= stockUpdate.Quantity
				json.NewEncoder(w).Encode(products[i])
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Insufficient stock"})
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products/{id}/stock", updateStockCount).Methods("PUT")

	log.Println("Product Service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
