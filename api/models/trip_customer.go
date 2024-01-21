package models

type TripCustomer struct {
	ID           string   `json:"id"`
	TripID       string   `json:"trip_id"`
	CustomerID   string   `json:"customer_id"`
	CustomerData Customer `json:"customer_data"`
	CreatedAt    string   `json:"created_at"`
}

type CreateTripCustomer struct {
	TripID     string `json:"trip_id"`
	CustomerID string `json:"customer_id"`
}

type TripCustomersResponse struct {
	TripCustomers []TripCustomer `json:"trip_customers"`
	Count         int            `json:"count"`
}