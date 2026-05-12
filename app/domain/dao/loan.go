package dao

type Loan struct {
	ID            int64   `bson:"id" json:"id"`
	UserID        string  `bson:"user_id" json:"user_id"`
	MRP           float64 `bson:"mrp" json:"mrp"`
	DP            float64 `bson:"dp" json:"dp"`
	VehicleYear   int     `bson:"vehicle_year" json:"vehicle_year"`
	PoliceNumber  string  `bson:"police_number" json:"police_number" gorm:"unique"`
	MachineNumber string  `bson:"machine_number" json:"machine_number"`
	Status        string  `bson:"status" json:"status"`
	RequestDate   string  `bson:"request_date" json:"request_date"`
	ApprovedDate  string  `bson:"approved_date" json:"approved_date"`
	CreatedAt     string  `bson:"created_at" json:"created_at"`
	UpdatedAt     string  `bson:"updated_at" json:"updated_at"`
}
