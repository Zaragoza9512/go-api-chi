#!/bin/bash

API_URL="http://localhost:8080"

echo "=== Testing API ==="
echo ""

# Login
echo "1. Login..."
TOKEN=$(curl -s -X POST $API_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' \
  | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "❌ Login failed"
  exit 1
fi
echo "✅ Login successful"
echo ""

# Get all products
echo "2. Getting all products..."
PRODUCTS=$(curl -s -X GET $API_URL/productos \
  -H "Authorization: Bearer $TOKEN")
echo "✅ Products retrieved"
echo "$PRODUCTS" | jq '.'
echo ""

# Create product
echo "3. Creating new product..."
NEW_PRODUCT=$(curl -s -X POST $API_URL/productos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Product",
    "description": "Testing API",
    "price": 99.99,
    "stock": 10
  }')
PRODUCT_ID=$(echo "$NEW_PRODUCT" | jq -r '.id')
echo "✅ Product created with ID: $PRODUCT_ID"
echo ""

# Get specific product
echo "4. Getting product $PRODUCT_ID..."
curl -s -X GET "$API_URL/productos/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.'
echo "✅ Product retrieved"
echo ""

# Update product
echo "5. Updating product $PRODUCT_ID..."
curl -s -X PUT "$API_URL/productos/$PRODUCT_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Product Updated",
    "description": "Updated description",
    "price": 149.99,
    "stock": 5
  }' | jq '.'
echo "✅ Product updated"
echo ""

# Delete product
echo "6. Deleting product $PRODUCT_ID..."
curl -s -X DELETE "$API_URL/productos/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN"
echo "✅ Product deleted"
echo ""

echo "=== All tests passed! ==="
