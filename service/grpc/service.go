package main

import (
	"context"
	"log"
	"net"

	pb "grpc/pb/gen/payment"

	"google.golang.org/grpc"
)

// server is used to implement payment.PaymentServiceServer.
type server struct {
	pb.UnimplementedPaymentServiceServer
}

// MakePayment implements payment.PaymentServiceServer
func (s *server) MakePayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	// Implement your payment logic here
	return &pb.PaymentResponse{
		Success:       true,
		TransactionId: "123456",
		Message:       "Payment successful",
		Timestamp:     "2023-10-01T12:00:00Z",
		Fee:           1.00,
	}, nil
}

// GetPaymentStatus implements payment.PaymentServiceServer
func (s *server) GetPaymentStatus(ctx context.Context, req *pb.PaymentStatusRequest) (*pb.PaymentStatusResponse, error) {
	// Implement your status retrieval logic here
	return &pb.PaymentStatusResponse{
		TransactionId: req.TransactionId,
		Status:        "completed",
		Message:       "Payment completed successfully",
		Timestamp:     "2023-10-01T12:00:00Z",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
