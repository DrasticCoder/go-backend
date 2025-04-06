#!/bin/bash

BASE_URL="http://localhost:8080/api/v1/auth"
PROFILE_URL="http://localhost:8080/api/v1/users/profile"

EMAIL="testuser@example.com"
PASSWORD="SecurePass123"
USERNAME="testuser"
ROLE="free"

echo "📝 Registering user..."
curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "'$USERNAME'",
    "email": "'$EMAIL'",
    "password": "'$PASSWORD'",
    "role": "'$ROLE'"
  }'
echo -e "\n✅ Registration complete.\n"

echo "🔐 Logging in..."
TOKEN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'$EMAIL'",
    "password": "'$PASSWORD'"
  }' | jq -r .token)

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Login failed. Check credentials or server logs."
  exit 1
fi

echo "✅ Login successful. JWT Token:"
echo "$TOKEN"
echo

echo "👤 Fetching user profile..."
curl -s -X GET "$PROFILE_URL" \
  -H "Authorization: Bearer $TOKEN" | jq
echo

echo "🚪 Logging out (no-op in stateless JWT)..."
curl -s -X POST "$BASE_URL/logout" \
  -H "Authorization: Bearer $TOKEN"
echo -e "\n✅ Logout done."

