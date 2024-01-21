package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type tripRepo struct {
	db *sql.DB
}

func NewTripRepo(db *sql.DB) storage.ITripRepo {
	return &tripRepo{
		db: db,
	}
}

func (c *tripRepo) Create(req models.CreateTrip) (string, error) {
	uid := uuid.New().String()

	if _, err := c.db.Exec(`insert into 
			trips (id, from_city_id, to_city_id, driver_id, price) values ($1, $2, $3, $4, $5)
			`,
		uid,
		req.FromCityID,
		req.ToCityID,
		req.DriverID,
		req.Price,
	); err != nil {
		fmt.Println("error while inserting data", err.Error())
		return "", err
	}

	return uid, nil

}

// TASK 9

func (c *tripRepo) Get(id string) (models.Trip, error) {

	trip := models.Trip{}

	query := `
		select id, trip_number_id, from_city_id, to_city_id, driver_id, price, created_at from trips where id = $1 
`
	if err := c.db.QueryRow(query, id).Scan(
		&trip.ID,
		&trip.TripNumberID,
		&trip.FromCityID,
		&trip.ToCityID,
		&trip.DriverID,
		&trip.Price,
		&trip.CreatedAt,
	); err != nil {
		fmt.Println("error while scanning trip", err.Error())
		return models.Trip{}, err
	}

	return trip, nil

}

// TASK 10

func (c *tripRepo) GetList(req models.GetListRequest) (models.TripsResponse, error) {

	var (
		trips             = []models.Trip{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
	)

	countQuery = `
		SELECT count(1) from trips `

	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of users", err.Error())
		return models.TripsResponse{}, err
	}

	query = `
    SELECT id, trip_number_id, from_city_id, to_city_id, driver_id, price, created_at
      FROM trips
          `

	query += ` LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.TripsResponse{}, err
	}

	for rows.Next() {
		trip := models.Trip{}

		if err = rows.Scan(
			&trip.ID,
			&trip.TripNumberID,
			&trip.FromCityID,
			&trip.ToCityID,
			&trip.DriverID,
			&trip.Price,
			&trip.CreatedAt,
		); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.TripsResponse{}, err
		}

		trips = append(trips, trip)
	}

	return models.TripsResponse{
		Trips: trips,
		Count: count,
	}, nil

}

func (c *tripRepo) Update(req models.Trip) (string, error) {
	query := `
		update trips 
			set from_city_id = $1, to_city_id = $2, driver_id = $3, price = $4
				where id = $5`

	if _, err := c.db.Exec(query, req.FromCityID, req.ToCityID, req.DriverID, req.Price, req.ID); err != nil {
		fmt.Println("error while updating trip data", err.Error())
		return "", err
	}

	return req.ID, nil

}

func (c *tripRepo) Delete(id string) error {
	query := `
		delete from trips
			where id = $1
`
	if _, err := c.db.Exec(query, id); err != nil {
		fmt.Println("error while deleting trip by id", err.Error())
		return err
	}

	return nil

}
