package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// struct untuk menyimpan data transaksi
type Transaction struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Customer  string    `json:"customer"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

// struct untuk menyimpan payload request dari user
type Request struct {
	RequestID uint          `json:"request_id"`
	Data      []Transaction `json:"data"`
}

// konfigurasi database
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "transaction_db"
)

// koneksi database
var db *gorm.DB
var err error

// menginisialisasi koneksi ke database
func initDB() {
	dbURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = gorm.Open("postgres", dbURI)
	if err != nil {
		panic("failed to connect database")
	}
	// migration
	db.AutoMigrate(&Transaction{})
}

// function untuk memproses request
func processRequest(request Request, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, data := range request.Data {
		db.Create(&data)
	}
}

// handler untuk endpoint /transaction dengan method POST
func addTransaction(w http.ResponseWriter, r *http.Request) {
	// dekode payload dari request
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// memproses request secara parallel menggunakan goroutine
	var wg sync.WaitGroup
	for i := 0; i < len(request.Data); i++ {
		wg.Add(1)
		go processRequest(request, &wg)
	}
	wg.Wait()

	// memberikan response ke user
	response := map[string]string{
		"message": "Data transaction berhasil dimasukkan",
	}
	json.NewEncoder(w).Encode(response)
}

// handler untuk endpoint /
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to transaction API!")
}

func NewRouter() {
	router := mux.NewRouter()

	// endpoint untuk /
	router.HandleFunc("/", home).Methods("GET")

	// endpoint untuk /transaction
	router.HandleFunc("/transaction", addTransaction).Methods("POST")
}

func main() {
	// inisialisasi koneksi ke database
	initDB()

	// setup router menggunakan mux
	router := mux.NewRouter()

	// endpoint untuk /
	router.HandleFunc("/", home).Methods("GET")

	// endpoint untuk /transaction
	router.HandleFunc("/transaction", addTransaction).Methods("POST")

	// jalankan server pada port 8080
	log.Fatal(http.ListenAndServe(":8080", router))
}
