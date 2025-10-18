package handlers

type sendOTPReq struct {
	Phone string `json:"phone" binding:"required"`
}

type verifyOTPReq struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

type onboardingReq struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email"`
	LicenseNumber   string `json:"license_number" binding:"required"`
	VehicleMake     string `json:"vehicle_make" binding:"required"`
	VehicleModel    string `json:"vehicle_model" binding:"required"`
	VehicleYear     int    `json:"vehicle_year" binding:"required"`
	PlateNumber     string `json:"plate_number" binding:"required"`
	InsuranceNumber string `json:"insurance_number" binding:"required"`
	RCNumber        string `json:"rc_number" binding:"required"`
}
