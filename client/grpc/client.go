package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "grpc/pb/gen/payment"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load("grpc_proto.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.NewClient(os.Getenv("GRPC_URL"), grpc.WithTransportCredentials((insecure.NewCredentials())))
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
