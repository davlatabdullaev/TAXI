package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"
	"time"

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
func (t tripRepo) Create(req models.CreateTrip) (string, error) {
	uid := uuid.New()
	createdAt := time.Now()

	tx, err := t.db.Begin()
	if err != nil {
		return "", fmt.Errorf("could not begin transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	if _, err := tx.Exec(`
		INSERT INTO trips (id, from_city_id, to_city_id, driver_id, price, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		`, uid, req.FromCityID, req.ToCityID, req.DriverID, req.Price, createdAt,
	); err != nil {
		tx.Rollback()
		return "", fmt.Errorf("error while inserting data: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("error committing transaction: %v", err)
	}

	return uid.String(), nil
}

func (c tripRepo) Get(id models.PrimaryKey) (models.Trip, error) {
	trip := models.Trip{}

	query := `
        SELECT
            t.id, 
			t.trip_number_id, 
			t.from_city_id, 
			t.to_city_id, 
			t.driver_id, 
			t.price, 
			t.created_at,
            cities_from.id AS from_city_data_id,
            cities_from.name AS from_city_data_name,
			cities_from.created_at AS from_city_data_created_at,
            cities_to.id AS to_city_data_id,
            cities_to.name AS to_city_data_name,
			cities_to.created_at AS to_city_data_created_at,
            drivers.id AS driver_data_id, 
			drivers.full_name AS driver_data_name,
			drivers.phone AS driver_data_phone,
			drivers.from_city_id AS driver_from_city_id,
			driver_from_cities.id AS driver_from_city_data_id,
			driver_from_cities.name AS driver_from_city_data_name,
			driver_from_cities.created_at AS driver_from_city_data_created_at,
			drivers.to_city_id AS driver_to_city_id,
			driver_to_cities.id AS driver_to_city_data_id,
			driver_to_cities.name AS driver_to_city_data_name,
			driver_to_cities.created_at AS driver_to_city_data_created_at,
			drivers.created_at AS driver_data_created_at
        FROM trips t
        JOIN cities cities_from ON t.from_city_id = cities_from.id
        JOIN cities cities_to ON t.to_city_id = cities_to.id
        JOIN drivers drivers ON t.driver_id = drivers.id
        JOIN cities driver_from_cities ON drivers.from_city_id = driver_from_cities.id
        JOIN cities driver_to_cities ON drivers.to_city_id = driver_to_cities.id
        WHERE t.id = $1
    `

	err := c.db.QueryRow(query, id.ID).Scan(
		&trip.ID,
		&trip.TripNumberID,
		&trip.FromCityID,
		&trip.ToCityID,
		&trip.DriverID,
		&trip.Price,
		&trip.CreatedAt,
		&trip.FromCityData.ID,
		&trip.FromCityData.Name,
		&trip.FromCityData.CreatedAt,
		&trip.ToCityData.ID,
		&trip.ToCityData.Name,
		&trip.ToCityData.CreatedAt,
		&trip.DriverData.ID,
		&trip.DriverData.FullName,
		&trip.DriverData.Phone,
		&trip.DriverData.FromCityID,
		&trip.DriverData.FromCityData.ID,
		&trip.DriverData.FromCityData.Name,
		&trip.DriverData.FromCityData.CreatedAt,
		&trip.DriverData.ToCityID,
		&trip.DriverData.ToCityData.ID,
		&trip.DriverData.ToCityData.Name,
		&trip.DriverData.ToCityData.CreatedAt,
		&trip.DriverData.CreatedAt,
	)

	if err != nil {
		fmt.Println("error while scanning trip and related data", err.Error())
		return models.Trip{}, err
	}

	return trip, nil
}

func (c tripRepo) GetList(req models.GetListRequest) (models.TripsResponse, error) {
	var (
		trips  = []models.Trip{}
		count  = 0
		page   = req.Page
		limit  = req.Limit
		offset = (page - 1) * limit
	)

	countQuery := `
        SELECT COUNT(1) FROM trips
    `

	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of trips", err.Error())
		return models.TripsResponse{}, err
	}

	query := `
        SELECT
            t.id, 
            t.trip_number_id, 
            t.from_city_id, 
            t.to_city_id, 
            t.driver_id, 
            t.price, 
            t.created_at,
            cities_from.id AS from_city_data_id,
            cities_from.name AS from_city_data_name,
            cities_from.created_at AS from_city_data_created_at,
            cities_to.id AS to_city_data_id,
            cities_to.name AS to_city_data_name,
            cities_to.created_at AS to_city_data_created_at,
            drivers.id AS driver_data_id, 
            drivers.full_name AS driver_data_name,
            drivers.phone AS driver_data_phone,
            drivers.from_city_id AS driver_from_city_id,
            driver_from_cities.id AS driver_from_city_data_id,
            driver_from_cities.name AS driver_from_city_data_name,
            driver_from_cities.created_at AS driver_from_city_data_created_at,
            drivers.to_city_id AS driver_to_city_id,
            driver_to_cities.id AS driver_to_city_data_id,
            driver_to_cities.name AS driver_to_city_data_name,
            driver_to_cities.created_at AS driver_to_city_data_created_at,
            drivers.created_at AS driver_data_created_at
        FROM trips t
        JOIN cities cities_from ON t.from_city_id = cities_from.id
        JOIN cities cities_to ON t.to_city_id = cities_to.id
        JOIN drivers drivers ON t.driver_id = drivers.id
        JOIN cities driver_from_cities ON drivers.from_city_id = driver_from_cities.id
        JOIN cities driver_to_cities ON drivers.to_city_id = driver_to_cities.id
        ORDER BY t.created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := c.db.Query(query, limit, offset)
	if err != nil {
		fmt.Println("error while querying rows", err.Error())
		return models.TripsResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		trip := models.Trip{}
		if err := rows.Scan(
			&trip.ID,
			&trip.TripNumberID,
			&trip.FromCityID,
			&trip.ToCityID,
			&trip.DriverID,
			&trip.Price,
			&trip.CreatedAt,
			&trip.FromCityData.ID,
			&trip.FromCityData.Name,
			&trip.FromCityData.CreatedAt,
			&trip.ToCityData.ID,
			&trip.ToCityData.Name,
			&trip.ToCityData.CreatedAt,
			&trip.DriverData.ID,
			&trip.DriverData.FullName,
			&trip.DriverData.Phone,
			&trip.DriverData.FromCityID,
			&trip.DriverData.FromCityData.ID,
			&trip.DriverData.FromCityData.Name,
			&trip.DriverData.FromCityData.CreatedAt,
			&trip.DriverData.ToCityID,
			&trip.DriverData.ToCityData.ID,
			&trip.DriverData.ToCityData.Name,
			&trip.DriverData.ToCityData.CreatedAt,
			&trip.DriverData.CreatedAt,
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

func (c tripRepo) Update(req models.Trip) (string, error) {
	query := `
        UPDATE trips 
        SET trip_number_id = $1, 
            from_city_id = $2, 
            to_city_id = $3, 
            driver_id = $4, 
            price = $5, 
            created_at = $6
        WHERE id = $7
    `

	_, err := c.db.Exec(query, req.TripNumberID, req.FromCityID, req.ToCityID, req.DriverID, req.Price, req.CreatedAt, req.ID)
	if err != nil {
		fmt.Println("error while updating trips data:", err.Error())
		return " ", err
	}

	return req.ID, nil
}

func (c tripRepo) Delete(id models.PrimaryKey) error {
	query := `
        delete from trips
        WHERE id = $1
    `
	if _, err := c.db.Exec(query, id.ID); err != nil {
		fmt.Println("error while deleting trip by id", err.Error())
		return err
	}

	return nil
}
