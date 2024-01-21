package handler

import (
	"city2city/api/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Car(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCar(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetCarList(w)
		} else {
			h.GetCarByID(w, r)
		}
	case http.MethodPut:
		values := r.URL.Query()
		if _, ok := values["route"]; ok {
			h.UpdateCarRoute(w, r)
		} else if _, ok := values["status"]; ok {
			h.UpdateCarStatus(w, r)
		} else {
			h.UpdateCar(w, r)
		}
	case http.MethodDelete:
		h.DeleteCar(w, r)
	}
}

func (h Handler) CreateCar(w http.ResponseWriter, r *http.Request) {
	createCar := models.CreateCar{}

	if err := json.NewDecoder(r.Body).Decode(&createCar); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Car().Create(createCar)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	car, err := h.storage.Car().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, car)

}

// TASK 3

func (h Handler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusInternalServerError, errors.New("id is required"))
		return
	}
	id := values["id"][0]
	var err error

	car, err := h.storage.Car().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, car)

}

// TASK 4

func (h Handler) GetCarList(w http.ResponseWriter) {

	var (
		page, limit = 1, 50
		err         error
	)

	response, err := h.storage.Car().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})

	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, response)

}

func (h Handler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	updateCar := models.Car{}

	if err := json.NewDecoder(r.Body).Decode(&updateCar); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Car().Update(updateCar)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	car, err := h.storage.Car().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, car)

}

func (h Handler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Car().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")

}

// TASK 1

func (h Handler) UpdateCarRoute(w http.ResponseWriter, r *http.Request) {
	updateCarRoute := models.UpdateCarRoute{}

	if err := json.NewDecoder(r.Body).Decode(&updateCarRoute); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(updateCarRoute.CarID) == 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	if err := h.storage.Car().UpdateCarRoute(updateCarRoute); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "Car route successfully updated")
}

// TASK 2

func (h Handler) UpdateCarStatus(w http.ResponseWriter, r *http.Request) {
	updateCarStatus := models.UpdateCarStatus{}

	if err := json.NewDecoder(r.Body).Decode(&updateCarStatus); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(updateCarStatus.ID) == 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	if updateCarStatus.ID == "" {
		handleResponse(w, http.StatusBadRequest, errors.New("id cannot be empty"))
		return
	}

	if err := h.storage.Car().UpdateCarStatus(updateCarStatus); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "car status updated successfully")
}
