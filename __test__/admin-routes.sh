# 🔐 Auth token
export ADMIN_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQwMTc2MzAsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOiI2N2YyM2I4NzAyYTQ3MzY2ZWZkYzNhZWIifQ.P6Nl80dpHCZU8yrLjW_LKzssoHV-ioOFG7hrCJemthk

# 🔐 1. List All Users (GET `/admin/users`)

curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"


# 🔍 2. Get Single User by ID (GET `/admin/users/:id`)

curl -X GET http://localhost:8080/api/v1/admin/users/USER_ID_HERE \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# ➕ 3. Create New User (POST `/admin/users`)

curl -X POST http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
        "username": "testuser",
        "email": "testuser@example.com",
        "password": "SecurePass123",
        "role": "free"
      }'

# ✏️ 4. Update User (PUT `/admin/users/:id`)

curl -X PUT http://localhost:8080/api/v1/admin/users/USER_ID_HERE \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
        "username": "updateduser",
        "role": "premium"
      }'


# ❌ 5. Delete User (DELETE `/admin/users/:id`)

curl -X DELETE http://localhost:8080/api/v1/admin/users/USER_ID_HERE \
  -H "Authorization: Bearer $ADMIN_TOKEN"


# 🧪 Bonus Tip — Quick test if role protection works:

# Try a route without a token or with a non-admin token to confirm RBAC kicks in:

curl -X GET http://localhost:8080/api/v1/admin/users

# 🔐 Forbidden: no role found in context
# 🔐 Forbidden: insufficient privileges
