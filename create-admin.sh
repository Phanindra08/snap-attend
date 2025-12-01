#!/bin/bash

# Script to create an admin user for SnapAttend
# This creates a Student user first, then you can use the admin panel to create more admins

echo "Creating initial admin user..."
echo ""
echo "Since the signup endpoint only allows Student/Professor,"
echo "we'll create a Student account that you can use to bootstrap."
echo ""
echo "Creating student account: admin@snapattend.com / admin12345"

curl -X POST http://localhost:8080/api/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Admin",
    "lastName": "User",
    "email": "admin@snapattend.com",
    "password": "admin12345",
    "profile": "Student",
    "address1": "123 Campus Dr",
    "address2": "",
    "city": "Charlotte",
    "state": "NC",
    "zipCode": "28262",
    "country": "USA"
  }'

echo ""
echo ""
echo "Account created! However, this is a STUDENT account."
echo "You need to manually update the database to make it an Admin."
echo ""
echo "Run this SQL command to convert the user to Admin:"
echo ""
echo "sqlite3 Backend/snapattend.db \"UPDATE users SET profile_id = (SELECT id FROM user_profiles WHERE role = 'Admin') WHERE email = 'admin@snapattend.com';\""
echo ""
