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
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPaymentServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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
		log.Fatalf("could not make payment: %v", err)
	}
	log.Printf("Payment Response: %v", paymentRes)

	// Get payment status
	statusReq := &pb.PaymentStatusRequest{
		TransactionId: paymentRes.TransactionId,
	}
	statusRes, err := c.GetPaymentStatus(ctx, statusReq)
	if err != nil {
		log.Fatalf("could not get payment status: %v", err)
	}
	log.Printf("Payment Status: %v", statusRes)
}
