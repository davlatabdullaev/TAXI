package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type driverRepo struct {
	DB *sql.DB
}

func NewDriverRepo(db *sql.DB) storage.IDriverRepo {
	return driverRepo{
		DB: db,
	}
}

func (d driverRepo) Create(driver models.CreateDriver) (string, error) {
	id := uuid.New()
	createdAt := time.Now()

	if _, err := d.DB.Exec(`INSERT INTO drivers VALUES ($1, $2, $3, $4, $5, $6)`,
		id, driver.FullName, driver.Phone, driver.FromCityID, driver.ToCityID, createdAt); err != nil {
		fmt.Println("error while inserting data", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (d driverRepo) Get(pkey models.PrimaryKey) (models.Driver, error) {
	driver := models.Driver{}
	if err := d.DB.QueryRow(`
        SELECT
            drivers.id,
            drivers.full_name,
            drivers.phone,
			drivers.from_city_id,
			cities_from.id AS from_city_data_id,
            cities_from.name AS from_city_data_name,
			cities_from.created_at AS from_city_data_created_at,
			drivers.to_city_id ,
			cities_to.id AS to_city_data_id,
            cities_to.name AS to_city_data_name,
			cities_to.created_at AS to_city_data_created_at,
			drivers.created_at
        FROM
            drivers
        LEFT JOIN
            cities AS cities_from ON drivers.from_city_id = cities_from.id
        LEFT JOIN
            cities AS cities_to ON drivers.to_city_id = cities_to.id
        WHERE
            drivers.id = $1
    `, pkey.ID).Scan(
		&driver.ID,
		&driver.FullName,
		&driver.Phone,
		&driver.FromCityID,
		&driver.FromCityData.ID,
		&driver.FromCityData.Name,
		&driver.FromCityData.CreatedAt,
		&driver.ToCityID,
		&driver.ToCityData.ID,
		&driver.ToCityData.Name,
		&driver.ToCityData.CreatedAt,
		&driver.CreatedAt,
	); err != nil {
		fmt.Println("error while querying driver by ID", err.Error())
		return models.Driver{}, err
	}

	return driver, nil
}

func (d driverRepo) GetList(request models.GetListRequest) (models.DriversResponse, error) {
	var (
		drivers = []models.Driver{}
		count   = 0
		query   string
	)

	countQuery := `
		SELECT count(1) FROM drivers
	`

	if err := d.DB.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of drivers", err.Error())
		return models.DriversResponse{}, err
	}

	query = `
		SELECT
			drivers.id,
			drivers.full_name,
			drivers.phone,
			drivers.from_city_id,
			cities_from.id AS from_city_data_id,
			cities_from.name AS from_city_data_name,
			cities_from.created_at AS from_city_data_created_at,
			drivers.to_city_id ,
			cities_to.id AS to_city_data_id,
			cities_to.name AS to_city_data_name,
			cities_to.created_at AS to_city_data_created_at,
			drivers.created_at
		FROM
			drivers
		LEFT JOIN
			cities AS cities_from ON drivers.from_city_id = cities_from.id
		LEFT JOIN
			cities AS cities_to ON drivers.to_city_id = cities_to.id
	`

	query += ` LIMIT $1 OFFSET $2`

	rows, err := d.DB.Query(query, request.Limit, (request.Page-1)*request.Limit)
	if err != nil {
		fmt.Println("error while querying rows", err.Error())
		return models.DriversResponse{}, err
	}

	for rows.Next() {
		var driver models.Driver

		if err = rows.Scan(
			&driver.ID,
			&driver.FullName,
			&driver.Phone,
			&driver.FromCityID,
			&driver.FromCityData.ID,
			&driver.FromCityData.Name,
			&driver.FromCityData.CreatedAt,
			&driver.ToCityID,
			&driver.ToCityData.ID,
			&driver.ToCityData.Name,
			&driver.ToCityData.CreatedAt,
			&driver.CreatedAt,
		); err != nil {
			fmt.Println("error while scanning row:", err)
			return models.DriversResponse{}, err
		}

		drivers = append(drivers, driver)
	}

	return models.DriversResponse{
		Drivers: drivers,
		Count:   count,
	}, nil
}

func (d driverRepo) Update(request models.Driver) (string, error) {

	query := `UPDATE drivers SET full_name = $1, phone = $2, from_city_id = $3, to_city_id = $4 WHERE id = $5`

	if _, err := d.DB.Exec(query, request.FullName, request.Phone, request.FromCityID, request.ToCityID, request.ID); err != nil {
		fmt.Println("error while updating driver data", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (d driverRepo) Delete(request models.PrimaryKey) error {

	query := `DELETE FROM drivers WHERE id = $1`

	if _, err := d.DB.Exec(query, request.ID); err != nil {
		fmt.Println("error while deleting driver by ID", err.Error())
		return err
	}

	return nil
}
