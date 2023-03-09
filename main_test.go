package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddTransactionHandler(t *testing.T) {
	// inisialisasi router dan server testing
	router := mux.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	// payload request
	transaction1 := Transaction{
		Customer:  "John Smith",
		Quantity:  2,
		Price:     10.50,
		Timestamp: time.Now(),
	}
	transaction2 := Transaction{
		Customer:  "Jane Doe",
		Quantity:  1,
		Price:     5.25,
		Timestamp: time.Now(),
	}
	requestData := Request{
		RequestID: 12345,
		Data:      []Transaction{transaction1, transaction2},
	}
	requestDataJSON, _ := json.Marshal(requestData)

	// membuat request ke endpoint /transaction
	req, err := http.NewRequest("POST", server.URL+"/transaction", bytes.NewBuffer(requestDataJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// melakukan request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// membaca response dari server
	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)

	// membandingkan response dengan yang diharapkan
	expectedResponse := map[string]string{
		"message": "Data transaction berhasil dimasukkan",
	}
	assert.Equal(t, expectedResponse, response)

	// memastikan data sudah dimasukkan ke database
	var count int
	db.Model(&Transaction{}).Count(&count)
	assert.Equal(t, 2, count)
}

func TestProcessRequest(t *testing.T) {
	// inisialisasi data transaksi dan wait group
	transaction1 := Transaction{
		Customer:  "John Smith",
		Quantity:  2,
		Price:     10.50,
		Timestamp: time.Now(),
	}
	transaction2 := Transaction{
		Customer:  "Jane Doe",
		Quantity:  1,
		Price:     5.25,
		Timestamp: time.Now(),
	}
	var wg sync.WaitGroup
	wg.Add(1)

	// memproses request secara parallel menggunakan goroutine
	requestData := Request{
		RequestID: 12345,
		Data:      []Transaction{transaction1, transaction2},
	}
	go processRequest(requestData, &wg)

	// menunggu proses selesai
	wg.Wait()

	// memastikan data sudah dimasukkan ke database
	var count int
	db.Model(&Transaction{}).Count(&count)
	assert.Equal(t, 2, count)
}
