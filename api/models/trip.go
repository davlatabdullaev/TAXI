package models

type Trip struct {
	ID           string `json:"id"`
	TripNumberID string `json:"trip_number_id"`
	FromCityID   string `json:"from_city_id"`
	FromCityData City   `json:"from_city_data"`
	ToCityID     string `json:"to_city_id"`
	ToCityData   City   `json:"to_city_data"`
	DriverID     string `json:"driver_id"`
	DriverData   Driver `json:"driver_data"`
	Price        int    `json:"price"`
	CreatedAt    string `json:"created_at"`
}

type CreateTrip struct {
	TripNumberID string `json:"trip_number_id"`
	FromCityID   string `json:"from_city_id"`
	ToCityID     string `json:"to_city_id"`
	DriverID     string `json:"driver_id"`
	Price        int    `json:"price"`
	CreatedAt    string `json:"created_at"`
}

type TripsResponse struct {
	Trips []Trip `json:"trips"`
	Count int    `json:"count"`
}