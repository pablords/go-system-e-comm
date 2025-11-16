#!/bin/bash

echo "🧪 Testing Orders ↔ Payments Integration"
echo ""

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if services are running
echo "📡 Checking if services are running..."

if ! curl -s http://localhost:8080/health > /dev/null; then
    echo -e "${RED}❌ Orders Service is not running on port 8080${NC}"
    echo "   Start it with: cd orders && go run cmd/api/main.go"
    exit 1
fi
echo -e "${GREEN}✅ Orders Service is running${NC}"

if ! nc -z localhost 50051 2>/dev/null; then
    echo -e "${RED}❌ Payment Service is not running on port 50051${NC}"
    echo "   Start it with: cd payments && ./bin/payment-service"
    exit 1
fi
echo -e "${GREEN}✅ Payment Service is running${NC}"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "SETUP: Creating sample products"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Create Product 1
PROD1_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product 1",
    "description": "Test product for integration",
    "price": 50.00
  }')
PRODUCT_ID_1=$(echo "$PROD1_RESPONSE" | jq -r '.id')
echo "✅ Product 1 created: $PRODUCT_ID_1"

# Create Product 2
PROD2_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product 2",
    "description": "Another test product",
    "price": 100.00
  }')
PRODUCT_ID_2=$(echo "$PROD2_RESPONSE" | jq -r '.id')
echo "✅ Product 2 created: $PRODUCT_ID_2"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "TEST 1: Create Order with Payment"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders/with-payment \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_email\": \"teste@example.com\",
    \"customer_name\": \"Cliente Teste\",
    \"payment_method\": 1,
    \"items\": [
      {
        \"product_id\": \"$PRODUCT_ID_1\",
        \"quantity\": 2,
        \"price\": 50.00
      },
      {
        \"product_id\": \"$PRODUCT_ID_2\",
        \"quantity\": 1,
        \"price\": 100.00
      }
    ]
  }")

echo "Response:"
echo "$RESPONSE" | jq '.'

ORDER_ID=$(echo "$RESPONSE" | jq -r '.order_id')
PAYMENT_ID=$(echo "$RESPONSE" | jq -r '.payment_id')
TOTAL=$(echo "$RESPONSE" | jq -r '.total')
STATUS=$(echo "$RESPONSE" | jq -r '.status')

if [ "$ORDER_ID" != "null" ] && [ "$ORDER_ID" != "" ]; then
    echo -e "${GREEN}✅ Order created successfully${NC}"
    echo "   Order ID: $ORDER_ID"
    echo "   Payment ID: $PAYMENT_ID"
    echo "   Total: R$ $TOTAL"
    echo "   Status: $STATUS"
else
    echo -e "${RED}❌ Failed to create order${NC}"
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "TEST 2: Get Order Details"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

ORDER_DETAILS=$(curl -s http://localhost:8080/api/v1/orders/$ORDER_ID)
echo "Order Details:"
echo "$ORDER_DETAILS" | jq '.'

if [ "$(echo "$ORDER_DETAILS" | jq -r '.id')" == "$ORDER_ID" ]; then
    echo -e "${GREEN}✅ Order retrieved successfully${NC}"
else
    echo -e "${RED}❌ Failed to retrieve order${NC}"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "TEST 3: Cancel Order and Payment"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CANCEL_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders/$ORDER_ID/cancel \
  -H "Content-Type: application/json" \
  -d "{
    \"payment_id\": \"$PAYMENT_ID\"
  }")

echo "Cancel Response:"
echo "$CANCEL_RESPONSE" | jq '.'

MESSAGE=$(echo "$CANCEL_RESPONSE" | jq -r '.message')
if [ "$MESSAGE" == "Order canceled successfully" ]; then
    echo -e "${GREEN}✅ Order and payment canceled successfully${NC}"
else
    echo -e "${YELLOW}⚠️  Unexpected cancel response${NC}"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "TEST 4: Verify Order Status After Cancel"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

ORDER_AFTER_CANCEL=$(curl -s http://localhost:8080/api/v1/orders/$ORDER_ID)
echo "Order After Cancel:"
echo "$ORDER_AFTER_CANCEL" | jq '.'

ORDER_STATUS=$(echo "$ORDER_AFTER_CANCEL" | jq -r '.status')
if [ "$ORDER_STATUS" == "canceled" ]; then
    echo -e "${GREEN}✅ Order status correctly updated to 'canceled'${NC}"
else
    echo -e "${YELLOW}⚠️  Order status is '$ORDER_STATUS', expected 'canceled'${NC}"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 Integration Tests Completed!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo ""
echo "🧹 Cleaning up test products..."
curl -s -X DELETE http://localhost:8080/api/v1/products/$PRODUCT_ID_1 > /dev/null
curl -s -X DELETE http://localhost:8080/api/v1/products/$PRODUCT_ID_2 > /dev/null
echo "✅ Test products deleted"
