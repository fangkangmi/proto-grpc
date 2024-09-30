package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func main() {
	ticker := time.NewTicker(200 * time.Millisecond) // 5 requests per second
	defer ticker.Stop()

	stop := time.After(1 * time.Hour)

	for {
		select {
		case <-stop:
			fmt.Println("Finished sending requests.")
			return
		case <-ticker.C:
			go func() {
				// Make a payment
				paymentReq := PaymentRequest{
					UserId:        "user123",
					Amount:        100.0,
					Currency:      "USD",
					PaymentMethod: "credit_card",
					Description:   "Payment for order #1234",
					RecipientId:   "recipient456",
				}
				paymentRes, err := makePayment(paymentReq)
				if err != nil {
					log.Printf("could not make payment: %v", err)
					return
				}
				fmt.Printf("Payment Response: %+v\n", paymentRes)

				// Get payment status
				statusReq := PaymentStatusRequest{
					TransactionId: paymentRes.TransactionId,
				}
				statusRes, err := getPaymentStatus(statusReq)
				if err != nil {
					log.Printf("could not get payment status: %v", err)
					return
				}
				fmt.Printf("Payment Status: %+v\n", statusRes)
			}()
		}
	}
}

func makePayment(req PaymentRequest) (*PaymentResponse, error) {
	url := "http://35.233.28.207/make_payment"
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func getPaymentStatus(req PaymentStatusRequest) (*PaymentStatusResponse, error) {
	url := "http://35.233.28.207/get_payment_status"
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res PaymentStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
