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
  echo "Error: REST_URL environment variable is not set."
  exit 1
fi

# Duration of the test
DURATION="15m"

# JSON payload for the POST request
PAYLOAD='{
  "user_id": "user123",
  "amount": 100.0,
  "currency": "USD",
  "payment_method": "credit_card",
  "description": "Payment for order #1234",
  "recipient_id": "recipient456"
}'

# Run the benchmark test without fixed RPS
hey -z $DURATION -m POST -H "Content-Type: application/json" -d "$PAYLOAD" $REST_URL