#!/bin/bash

# Source the .env file to set environment variables
if [ -f .env ]; then
  source .env
else
  echo ".env file not found."
  exit 1
fi

# REST service URL
if [ -z "$REST_URL" ]; then
  echo "Error: GRPC_URL environment variable is not set."
  exit 1
fi

# Number of requests per second
RPS=500

# Duration of the test
DURATION="15m"

# Function to generate a 2KB string without newline characters
generate_2kb_data() {
  local size=2048
  head -c $size </dev/urandom | base64 | tr -d '\n' | head -c $size
}

# Generate 2KB payload
METADATA=$(generate_2kb_data)

# JSON payload for the POST request
PAYLOAD=$(cat <<EOF
{
  "user_id": "user123",
  "amount": 100.0,
  "currency": "USD",
  "payment_method": "credit_card",
  "description": "$METADATA",
  "recipient_id": "recipient456"
}
EOF
)

# Run the load test
hey -t 1 -o "csv" -z $DURATION -q $RPS -m POST -H "Content-Type: application/json" -d "$PAYLOAD" $REST_URL
