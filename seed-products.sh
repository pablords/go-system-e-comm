#!/bin/bash

echo "🛍️  Creating sample products in Orders database..."

# Database credentials
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-root}
DB_PASSWORD=${DB_PASSWORD:-root}
DB_NAME=${DB_NAME:-orders_db}

# Check if mysql is available
# if ! command -v mysql &> /dev/null; then
#     echo "❌ mysql command not found. Please install mysql client."
#     exit 1
# fi

# Create products using HTTP API
echo ""
echo "Creating products via API..."

# Product 1
echo "📦 Creating Product 1: Notebook Dell"
RESPONSE1=$(curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Notebook Dell Inspiron 15",
    "description": "Notebook Dell Inspiron 15 com Intel Core i7, 16GB RAM, 512GB SSD",
    "price": 3500.00,
    "stock": 10
  }')

PRODUCT_1=$(echo "$RESPONSE1" | jq -r '.id')
echo "   ✅ Created: $PRODUCT_1"

# Product 2
echo "📦 Creating Product 2: Mouse Logitech"
RESPONSE2=$(curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mouse Logitech MX Master 3",
    "description": "Mouse sem fio ergonômico Logitech MX Master 3",
    "price": 450.00,
    "stock": 25
  }')

PRODUCT_2=$(echo "$RESPONSE2" | jq -r '.id')
echo "   ✅ Created: $PRODUCT_2"

# Product 3
echo "📦 Creating Product 3: Teclado Mecânico"
RESPONSE3=$(curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Teclado Mecânico Keychron K2",
    "description": "Teclado mecânico sem fio Keychron K2 com switch Gateron Brown",
    "price": 650.00,
    "stock": 15
  }')

PRODUCT_3=$(echo "$RESPONSE3" | jq -r '.id')
echo "   ✅ Created: $PRODUCT_3"

# Product 4
echo "📦 Creating Product 4: Monitor LG"
RESPONSE4=$(curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Monitor LG UltraWide 34\"",
    "description": "Monitor LG UltraWide 34\" 21:9 IPS 144Hz",
    "price": 2200.00,
    "stock": 8
  }')

PRODUCT_4=$(echo "$RESPONSE4" | jq -r '.id')
echo "   ✅ Created: $PRODUCT_4"

# Product 5
echo "📦 Creating Product 5: Webcam Logitech"
RESPONSE5=$(curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Webcam Logitech C920",
    "description": "Webcam Full HD 1080p Logitech C920 com microfone estéreo",
    "price": 380.00,
    "stock": 20
  }')

PRODUCT_5=$(echo "$RESPONSE5" | jq -r '.id')
echo "   ✅ Created: $PRODUCT_5"

# Product 6
echo "📦 Creating Product 6: Fone Bluetooth"
RESPONSE6=$(curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Fone Sony WH-1000XM5",
    "description": "Fone de ouvido Bluetooth Sony WH-1000XM5 com cancelamento de ruído",
    "price": 1800.00,
    "stock": 12
  }')

PRODUCT_6=$(echo "$RESPONSE6" | jq -r '.id')
echo "   ✅ Created: $PRODUCT_6"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Products created successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📋 Product IDs for testing:"
echo "   1. Notebook Dell: $PRODUCT_1"
echo "   2. Mouse Logitech: $PRODUCT_2"
echo "   3. Teclado Mecânico: $PRODUCT_3"
echo "   4. Monitor LG: $PRODUCT_4"
echo "   5. Webcam Logitech: $PRODUCT_5"
echo "   6. Fone Sony: $PRODUCT_6"
echo ""
echo "You can now use these product IDs to create orders!"
echo ""
echo "Example:"
echo "curl -X POST http://localhost:8080/api/v1/orders/with-payment \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{"
echo "    \"customer_email\": \"cliente@example.com\","
echo "    \"customer_name\": \"João Silva\","
echo "    \"payment_method\": 1,"
echo "    \"items\": ["
echo "      {\"product_id\": \"'$PRODUCT_2'\", \"quantity\": 1, \"price\": 450.00},"
echo "      {\"product_id\": \"'$PRODUCT_3'\", \"quantity\": 1, \"price\": 650.00}"
echo "    ]"
echo "  }'"
