package main

import (
	"context"
	"log"
	"time"

	pb "grpc/pb/gen/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.NewClient("34.38.55.203:80", grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPaymentServiceClient(conn)

	// Create a ticker that ticks every 200 milliseconds (5 times per second)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	// Create a channel that will send a signal after 1 hour
	stop := time.After(1 * time.Hour)

	for {
		select {
		case <-stop:
			log.Println("Finished sending requests.")
			return
		case <-ticker.C:
			go func() {
				// Contact the server and print out its response.
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				// Make a payment
				paymentReq := &pb.PaymentRequest{
					UserId:        "user123",
					Amount:        100.0,
					Currency:      "USD",
					PaymentMethod: "credit_card",
					Description:   "Payment for order #1234",
					RecipientId:   "recipient456",
				}
				paymentRes, err := c.MakePayment(ctx, paymentReq)
				if err != nil {
					log.Printf("could not make payment: %v", err)
					return
				}
				log.Printf("Payment Response: %v", paymentRes)

				// Get payment status
				statusReq := &pb.PaymentStatusRequest{
					TransactionId: paymentRes.TransactionId,
				}
				statusRes, err := c.GetPaymentStatus(ctx, statusReq)
				if err != nil {
					log.Printf("could not get payment status: %v", err)
					return
				}
				log.Printf("Payment Status: %v", statusRes)
			}()
		}
	}
}
