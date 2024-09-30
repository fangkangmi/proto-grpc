#!/bin/bash

# Source the .env file to set environment variables
if [ -f .env ]; then
  source .env
else
  echo ".env file not found."
  exit 1
fi

# gRPC service URL
if [ -z "$GRPC_URL" ]; then
  echo "Error: GRPC_URL environment variable is not set."
  exit 1
fi

# Number of requests per second
RPS=500

# Duration of the test (in minutes)
DURATION="15m"

# Total number of requests (RPS * 60 seconds * 15 minutes)
TOTAL_REQUESTS=$((RPS * 60 * 15))

# Function to generate a 2KB string without newline characters
generate_2kb_data() {
  local size=2048
  head -c $size </dev/urandom | base64 | tr -d '\n' | head -c $size
}

# Generate 2KB payload
PAYLOAD=$(generate_2kb_data)

# Run the load test
ghz --insecure \
    --proto ./pb/payment.proto \
    --call payment.PaymentService.MakePayment \
    --rps $RPS \
    --duration $DURATION \
    --data "{\"user_id\":\"user123\",\"amount\":100.0,\"currency\":\"USD\",\"payment_method\":\"credit_card\",\"description\":\"$PAYLOAD\",\"recipient_id\":\"recipient456\"}" \
    $GRPC_URL