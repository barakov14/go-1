package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const (
	redisAddress      = "localhost:6379"
	sqliteDatabaseURL = "./products.db"
)

var (
	db    *sql.DB
	cache *redis.Client
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func init() {
	// Initialize Redis client
	cache = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Initialize SQLite database
	var err error
	db, err = sql.Open("sqlite3", sqliteDatabaseURL)
	if err != nil {
		log.Fatal("Error connecting to the SQLite database:", err)
	}
}

func getProductFromDB(id int) (Product, error) {
	var product Product
	row := db.QueryRow("SELECT id, name, description, price FROM products WHERE id = ?", id)
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	return product, err
}

func getProductFromCache(id int) (Product, error) {
	var product Product
	val, err := cache.Get(context.Background(), strconv.Itoa(id)).Result()
	if err == redis.Nil {
		return product, sql.ErrNoRows
	} else if err != nil {
		return product, err
	}
	err = json.Unmarshal([]byte(val), &product)
	return product, err
}

func setProductToCache(id int, product Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	err = cache.Set(context.Background(), strconv.Itoa(id), data, 24*time.Hour).Err()
	return err
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Check if product is in cache
	product, err := getProductFromCache(id)
	if err == nil {
		json.NewEncoder(w).Encode(product)
		return
	}

	// If not in cache, fetch from database
	product, err = getProductFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Set product to cache
	err = setProductToCache(id, product)
	if err != nil {
		log.Println("Failed to set product to cache:", err)
	}

	json.NewEncoder(w).Encode(product)
}

func main() {
	db, err := sql.Open("sqlite3", sqliteDatabaseURL)
	if err != nil {
		log.Fatal("Error connecting to the SQLite database:", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY,
		name TEXT,
		description TEXT,
		price REAL
	)`)
	if err != nil {
		log.Fatal("Error creating products table:", err)
	}

	// Insert test data
	products := []struct {
		ID          int
		Name        string
		Description string
		Price       float64
	}{
		{1, "Product 1", "Description of Product 1", 10.99},
		{2, "Product 2", "Description of Product 2", 20.99},
		{3, "Product 3", "Description of Product 3", 30.99},
	}

	for _, p := range products {
		_, err := db.Exec("INSERT INTO products (id, name, description, price) VALUES (?, ?, ?, ?)", p.ID, p.Name, p.Description, p.Price)
		if err != nil {
			log.Fatal("Error inserting test data:", err)
		}
	}

	log.Println("Test data inserted successfully")

	r := mux.NewRouter()
	r.HandleFunc("/products/{id}", getProductHandler).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
