package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type PaymentRequest struct {
	UserId        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
	Description   string  `json:"description"`
	RecipientId   string  `json:"recipient_id"`
}

type PaymentResponse struct {
	Success       bool    `json:"success"`
	TransactionId string  `json:"transaction_id"`
	Message       string  `json:"message"`
	Timestamp     string  `json:"timestamp"`
	Fee           float64 `json:"fee"`
}

type PaymentStatusRequest struct {
	TransactionId string `json:"transaction_id"`
}

type PaymentStatusResponse struct {
	TransactionId string `json:"transaction_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	Timestamp     string `json:"timestamp"`
}

var (
	payments = make(map[string]PaymentResponse)
	mu       sync.Mutex
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/make_payment", makePaymentHandler)
	mux.HandleFunc("/get_payment_status", getPaymentStatusHandler)

	server := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on :8081")
	log.Fatal(server.ListenAndServe())
}

func makePaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simulate payment processing
	transactionId := "txn_" + time.Now().Format("20060102150405")
	timestamp := time.Now().Format(time.RFC3339)
	res := PaymentResponse{
		Success:       true,
		TransactionId: transactionId,
		Message:       "Payment successful",
		Timestamp:     timestamp,
		Fee:           1.00,
	}

	// Store the payment response
	mu.Lock()
	payments[transactionId] = res
	mu.Unlock()

	json.NewEncoder(w).Encode(res)
}

func getPaymentStatusHandler(w http.ResponseWriter, r *http.Request) {
	var req PaymentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the payment status
	mu.Lock()
	res, exists := payments[req.TransactionId]
	mu.Unlock()

	if !exists {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	statusRes := PaymentStatusResponse{
		TransactionId: res.TransactionId,
		Status:        "completed",
		Message:       res.Message,
		Timestamp:     res.Timestamp,
	}
	json.NewEncoder(w).Encode(statusRes)
}
