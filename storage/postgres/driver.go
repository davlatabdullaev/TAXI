package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"

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
	uid := uuid.New()
	if _, err := d.DB.Exec(`insert into 
			drivers values ($1, $2, $3, $4, $5)
			`,
		uid,
		driver.FullName,
		driver.Phone,
		driver.FromCityID,
		driver.ToCityID,
	); err != nil {
		fmt.Println("error while inserting data", err.Error())
		return "", err
	}

	return uid.String(), nil

}

// TASK 5

func (d driverRepo) Get(id string) (models.Driver, error) {

	driver := models.Driver{}

	query := `
		select id, full_name, phone, from_city_id, to_city_id, created_at from drivers where id = $1
`
	if err := d.DB.QueryRow(query, id).Scan(
		&driver.ID,
		&driver.FullName,
		&driver.Phone,
		&driver.FromCityID,
		&driver.ToCityID,
		&driver.CreatedAt,
	); err != nil {
		fmt.Println("error while scanning user", err.Error())
		return models.Driver{}, err
	}

	return driver, nil

}

// TASK 6

func (d driverRepo) GetList(req models.GetListRequest) (models.DriversResponse, error) {
	var (
		drivers           = []models.Driver{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
	)

	countQuery = `
		SELECT count(1) from drivers`

	if err := d.DB.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of drivers", err.Error())
		return models.DriversResponse{}, err
	}

	query = `
		SELECT id, full_name, phone, from_city_id, to_city_id, created_at
			FROM drivers
			    `

	query += ` LIMIT $1 OFFSET $2`

	rows, err := d.DB.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.DriversResponse{}, err
	}

	for rows.Next() {
		driver := models.Driver{}

		if err = rows.Scan(
			&driver.ID,
			&driver.FullName,
			&driver.Phone,
			&driver.FromCityID,
			&driver.ToCityID,
			&driver.CreatedAt,
		); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.DriversResponse{}, err
		}

		drivers = append(drivers, driver)
	}

	return models.DriversResponse{
		Drivers: drivers,
		Count:   count,
	}, nil

}

func (d driverRepo) Update(driver models.Driver) (string, error) {

	query := `
	update drivers 
		set full_name = $1, phone = $2, from_city_id = $3, to_city_id = $4
			where id = $5`

	if _, err := d.DB.Exec(query, driver.FullName, driver.Phone, driver.FromCityID, driver.ToCityID, driver.ID); err != nil {
		fmt.Println("error while updating driver data", err.Error())
		return "", err
	}
	return driver.ID, nil

}

func (d driverRepo) Delete(id string) error {
	query := `
	delete from drivers
		where id = $1
`
	if _, err := d.DB.Exec(query, id); err != nil {
		fmt.Println("error while deleting driver by id", err.Error())
		return err
	}

	return nil
}
