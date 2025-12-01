#!/bin/bash
API="http://localhost:8080/api"

# Login as Admin
ADMIN_TOKEN=$(curl -s -X POST $API/auth/login -H "Content-Type: application/json" -d '{"email":"admin@snapattend.com", "password":"admin12345", "profile":"Admin"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

# Create Student
STUDENT_ID=$(curl -s -X POST $API/admin/students -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d '{"firstName":"Auth","lastName":"Test","email":"auth_test_unique@test.com","password":"password123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['student']['id'])")

echo "Created Student ID: $STUDENT_ID"
