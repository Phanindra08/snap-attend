#!/bin/bash
API="http://localhost:8080/api"

# Admin Login
ADMIN_TOKEN=$(curl -s -X POST $API/auth/login -H "Content-Type: application/json" -d '{"email":"admin@snapattend.com", "password":"admin12345", "profile":"Admin"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

# Create Course
COURSE_ID=$(curl -s -X POST $API/admin/courses -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d '{"courseName":"Question Flow Test"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['course']['id'])")

# Create Section
SECTION_ID=$(curl -s -X POST $API/admin/sections -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d "{\"courseId\":$COURSE_ID, \"sectionNumber\":\"Q101\", \"semesterId\":1, \"roomId\":1, \"latitude\":40.7128, \"longitude\":-74.0060}" | python3 -c "import sys, json; print(json.load(sys.stdin)['section']['id'])")

# Create Professor
PROF_EMAIL="prof_q@test.com"
PROF_PASS="password123"
PROF_ID=$(curl -s -X POST $API/admin/professors -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d "{\"firstName\":\"Prof\",\"lastName\":\"Q\",\"email\":\"$PROF_EMAIL\",\"password\":\"$PROF_PASS\"}" | python3 -c "import sys, json; print(json.load(sys.stdin)['professor']['id'])")

# Assign Professor
curl -s -X POST $API/admin/sections/assign-professor -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d "{\"professorId\":$PROF_ID, \"sectionId\":$SECTION_ID}"

# Professor Login
PROF_TOKEN=$(curl -s -X POST $API/auth/login -H "Content-Type: application/json" -d "{\"email\":\"$PROF_EMAIL\", \"password\":\"$PROF_PASS\", \"profile\":\"Professor\"}" | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

# Generate QR
QR_ID=$(curl -s -X POST $API/professor/qr -H "Authorization: Bearer $PROF_TOKEN" -H "Content-Type: application/json" -d "{\"sectionId\":$SECTION_ID, \"validSeconds\":3600}" | python3 -c "import sys, json; print(json.load(sys.stdin)['qr']['id'])")

echo "SECTION_ID=$SECTION_ID"
echo "QR_ID=$QR_ID"
echo "PROF_EMAIL=$PROF_EMAIL"
echo "PROF_PASS=$PROF_PASS"
