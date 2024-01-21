package models

type Customer struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type CreateCustomer struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type CustomersResponse struct {
	Customers []Customer `json:"customers"`
	Count     int        `json:"count"`
}