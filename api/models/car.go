package models

type Car struct {
	ID         string `json:"id"`
	Model      string `json:"model"`
	Brand      string `json:"brand"`
	Number     string `json:"number"`
	Status     string `json:"status"`
	DriverID   string `json:"driver_id"`
	DriverData Driver `json:"driver_data"`
	CreatedAt  string `json:"created_at"`
}

type CreateCar struct {
	Model    string `json:"model"`
	Brand    string `json:"brand"`
	Number   string `json:"number"`
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
