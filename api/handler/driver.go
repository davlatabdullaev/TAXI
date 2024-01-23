package handler

import (
	"city2city/api/models"
	"city2city/check"
	"encoding/json"
	"errors"
	"net/http"
)

func (h Handler) Driver(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateDriver(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		_, ok := values["id"]
		if !ok {
			h.GetDriverList(w, r)
		} else {
			h.GetDriverByID(w, r)
		}
	case http.MethodPut:
		h.UpdateDriver(w, r)
	case http.MethodDelete:
		h.DeleteDriver(w, r)
	}
}

func (h Handler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	createDriver := models.CreateDriver{}

	if err := json.NewDecoder(r.Body).Decode(&createDriver); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	if !check.PhoneNumber(createDriver.Phone) {
		handleResponse(w, http.StatusBadRequest, errors.New("phone number is not correct!"))
		return
	}

	id, err := h.storage.Driver().Create(createDriver)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	driver, err := h.storage.Driver().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, driver)

}

// TASK 5

func (h Handler) GetDriverByID(w http.ResponseWriter, r *http.Request) {

	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]
	var err error

	user, err := h.storage.Driver().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, user)

}

// TASK 6

func (h Handler) GetDriverList(w http.ResponseWriter, r *http.Request) {
	var (
		page, limit = 1, 50
		err         error
	)

	resp, err := h.storage.Driver().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, resp)
}

func (h Handler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	updateDriver := models.Driver{}

	if err := json.NewDecoder(r.Body).Decode(&updateDriver); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Driver().Update(updateDriver)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	driver, err := h.storage.Driver().Get(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, driver)

}

func (h Handler) DeleteDriver(w http.ResponseWriter, r *http.Request) {

	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	err := h.storage.Driver().Delete(models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data successfully deleted")

}
