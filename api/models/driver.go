package models

type Driver struct {
	ID           string `json:"id"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	FromCityID   string `json:"from_city_id"`
	FromCityData City   `json:"from_city_data"`
	ToCityID     string `json:"to_city_id"`
	ToCityData   City   `json:"to_city_data"`
	CreatedAt    string `json:"created_at"`
}

type CreateDriver struct {
	FullName   string `json:"full_name"`
	Phone      string `json:"phone"`
	FromCityID string `json:"from_city_id"`
	ToCityID   string `json:"to_city_id"`
}

type DriversResponse struct {
	Drivers []Driver `json:"drivers"`
	Count   int      `json:"count"`
}