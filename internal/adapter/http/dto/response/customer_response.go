package response

// CustomerVehicleResponse represents a vehicle associated with a customer.
type CustomerVehicleResponse struct {
	ID    uint   `json:"id"`
	Brand string `json:"brand,omitempty"`
	Model string `json:"model,omitempty"`
	Year  string `json:"year,omitempty"`
	Plate string `json:"plate,omitempty"`
}

// CustomerResponse represents the payload returned for customer operations.
type CustomerResponse struct {
	ID          uint                      `json:"id"`
	UserID      uint                      `json:"user_id,omitempty"`
	FullName    string                    `json:"full_name"`
	Email       string                    `json:"email"`
	PhoneNumber string                    `json:"phone_number"`
	Document    string                    `json:"document"`
	Vehicles    []CustomerVehicleResponse `json:"vehicles,omitempty"`
}
