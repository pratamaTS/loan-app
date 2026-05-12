package dto

type RequestLoanRequest struct {
	UserID        string  `json:"user_id" binding:"required"`
	MRP           float64 `json:"mrp" binding:"required"`
	DP            float64 `json:"dp" binding:"required"`
	VehicleYear   int     `json:"vehicle_year" binding:"required"`
	PoliceNumber  string  `json:"police_number" binding:"required"`
	MachineNumber string  `json:"machine_number" binding:"required"`
}

type ApproveLoanRequest struct {
	UserID       string `json:"user_id" binding:"required"`
	PoliceNumber string `json:"police_number" binding:"required"`
}
