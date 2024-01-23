package api

import (
	"city2city/api/handler"
	"net/http"
)

func New(h handler.Handler) {

	http.HandleFunc("/city", h.City)
	http.HandleFunc("/customer", h.Customer)
	http.HandleFunc("/driver", h.Driver)
	http.HandleFunc("/car", h.Car)
	http.HandleFunc("/trip", h.Trip)
	http.HandleFunc("/trip_customer", h.TripCustomer)
}