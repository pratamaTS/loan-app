Loan App — Loan Application Service

- Overview
Service sederhana untuk pengajuan dan persetujuan loan.

- How to Run
go run main.go
Service akan berjalan di:
http://localhost:10110

- Health Check
curl -X GET http://localhost:10110/health

- Request Loan
curl -X POST http://localhost:10110/api/loan/request \  -H "Content-Type: application/json" \  -d '{    "user_id": "Bruce",    "mrp": 10000000,    "dp": 2000000,    "vehicle_year": 2022,    "police_number": "B 1234 BYE",    "machine_number": "ENG12345"  }'

- Approve Loan
curl -X POST http://localhost:10110/api/loan/approve \  -H "Content-Type: application/json" \  -d '{    "user_id": "Bruce",    "police_number": "B 1234 BYE"  }'
