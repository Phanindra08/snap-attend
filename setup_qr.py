import requests
import json
import sys

API = "http://localhost:8080/api"

def get_token(email, password, profile):
    res = requests.post(f"{API}/auth/login", json={"email": email, "password": password, "profile": profile})
    res.raise_for_status()
    return res.json()["token"]

try:
    # Admin Login
    admin_token = get_token("admin@snapattend.com", "admin12345", "Admin")
    print(f"Admin Token obtained")

    # Assign Prof 1 to Section 10
    headers = {"Authorization": f"Bearer {admin_token}"}
    res = requests.post(f"{API}/admin/sections/assign-professor", json={"professorId": 1, "sectionId": 10}, headers=headers)
    # Ignore error if already assigned
    print("Assigned Professor")

    # Prof Login
    prof_token = get_token("prof@snapattend.com", "password123", "Professor")
    print(f"Prof Token obtained")

    # Generate QR
    prof_headers = {"Authorization": f"Bearer {prof_token}"}
    res = requests.post(f"{API}/professor/qr", json={"sectionId": 10, "validSeconds": 3600}, headers=prof_headers)
    res.raise_for_status()
    qr_id = res.json()["qr"]["id"]
    print(f"QR_ID={qr_id}")
    print(f"SECTION_ID=10")

except Exception as e:
    print(f"Error: {e}")
    sys.exit(1)
