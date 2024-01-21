package models

import "time"

type Car struct {
	ID         string `json:"id"`
	Model      string `json:"model"`
	Brand      string `json:"brand"`
	Number     string `json:"number"` // Number type int dan string ga o'zgartirildi
	Status     string `json:"status"`
	DriverID   string `json:"driver_id"`
	DriverData Driver `json:"driver_data"`
	CreatedAt  string `json:"created_at"`
}

type CreateCar struct {
	Model    string `json:"model"`
	Brand    string `json:"brand"`
	Number   string `json:"number"` // Number type int dan string ga o'zgartirildi
	DriverID string `json:"driver_id"`
}

type CarsResponse struct {
	Cars  []Car `json:"cars"`
	Count int   `json:"count"`
}

type UpdateCarStatus struct {
	ID     string `json:"id"`
	Status bool   `json:"status"`
}

type UpdateCarRoute struct {
	CarID         string    `json:"car_id"`
	DepartureTime time.Time `json:"departure_time"`
	FromCityID    string    `json:"from_city_id"`
	ToCityID      string    `json:"to_city_id"`
}
