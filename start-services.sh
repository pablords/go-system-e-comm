#!/bin/bash

echo "🚀 Starting Payment Service..."
cd payments
./bin/payments-service &
PAYMENT_PID=$!
echo "✅ Payment Service started (PID: $PAYMENT_PID)"

echo ""
echo "⏳ Waiting 3 seconds for Payment Service to initialize..."
sleep 3

echo ""
echo "🚀 Starting Orders Service..."
cd ../orders
go run cmd/api/main.go &
ORDERS_PID=$!
echo "✅ Orders Service started (PID: $ORDERS_PID)"

echo ""
echo "🎉 Both services are running!"
echo "   Payment Service: localhost:50051 (gRPC)"
echo "   Orders Service:  localhost:8080 (HTTP)"
echo ""
echo "Press Ctrl+C to stop all services..."

# Wait for Ctrl+C
trap "kill $PAYMENT_PID $ORDERS_PID; exit" INT TERM

wait
