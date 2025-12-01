#!/bin/bash

# Base URL
API="http://localhost:8080/api"

# 1. Login as Admin
echo "Logging in as Admin..."
ADMIN_TOKEN=$(curl -s -X POST $API/auth/login -H "Content-Type: application/json" -d '{"email":"admin@snapattend.com", "password":"admin12345", "profile":"Admin"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
echo "Admin Token: $ADMIN_TOKEN"

# 2. Create Course
echo "Creating Course..."
COURSE_ID=$(curl -s -X POST $API/admin/courses -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d '{"courseName":"E2E Test Course"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['course']['id'])")
echo "Course ID: $COURSE_ID"

# 3. Create Section with Coordinates (NYC)
echo "Creating Section..."
SECTION_ID=$(curl -s -X POST $API/admin/sections -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d "{\"courseId\":$COURSE_ID, \"sectionNumber\":\"999\", \"semesterId\":1, \"roomId\":1, \"latitude\":40.7128, \"longitude\":-74.0060}" | python3 -c "import sys, json; print(json.load(sys.stdin)['section']['id'])")
echo "Section ID: $SECTION_ID"

# 4. Login as Professor (using existing one or creating new? Let's use existing prof@snapattend.com if exists, or create one)
# We'll assume prof@snapattend.com exists from seed or previous tests.
echo "Logging in as Professor..."
PROF_TOKEN=$(curl -s -X POST $API/auth/login -H "Content-Type: application/json" -d '{"email":"prof@snapattend.com", "password":"password123", "profile":"Professor"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

if [ -z "$PROF_TOKEN" ]; then
    echo "Creating Professor..."
    PROF_ID=$(curl -s -X POST $API/admin/professors -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d '{"firstName":"Test","lastName":"Prof","email":"e2e_prof@test.com","password":"password123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['professor']['id'])")
    PROF_TOKEN=$(curl -s -X POST $API/auth/login -H "Content-Type: application/json" -d '{"email":"e2e_prof@test.com", "password":"password123", "profile":"Professor"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
    # Assign prof to section
    curl -s -X POST $API/admin/sections/assign-professor -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d "{\"professorId\":$PROF_ID, \"sectionId\":$SECTION_ID}"
else
    # Assign existing prof (ID 1 usually)
    curl -s -X POST $API/admin/sections/assign-professor -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d "{\"professorId\":1, \"sectionId\":$SECTION_ID}"
fi

# 5. Generate QR Code
echo "Generating QR Code..."
QR_ID=$(curl -s -X POST $API/professor/qr -H "Authorization: Bearer $PROF_TOKEN" -H "Content-Type: application/json" -d "{\"sectionId\":$SECTION_ID, \"validSeconds\":3600}" | python3 -c "import sys, json; print(json.load(sys.stdin)['qr']['id'])")
echo "QR ID: $QR_ID"

echo "SETUP_COMPLETE"
echo "QR_ID=$QR_ID"
echo "SECTION_ID=$SECTION_ID"
