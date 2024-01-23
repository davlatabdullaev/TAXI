package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type carRepo struct {
	db *sql.DB
}

func NewCarRepo(db *sql.DB) storage.ICarRepo {
	return carRepo{
		db,
	}
}

func (c carRepo) Create(car models.CreateCar) (string, error) {

	uid := uuid.New().String()
	query := `INSERT INTO cars (id, model, brand, number, driver_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := c.db.Exec(query, uid, car.Model, car.Brand, car.Number, car.DriverID)
	if err != nil {
		fmt.Println("error while inserting data ", err.Error())
		return "", err
	}

	return uid, nil

}

// TASK 3

func (c carRepo) Get(id string) (models.Car, error) {
	car := models.Car{}
	query := `
		SELECT
			c.id,
			c.model,
			c.brand,
			c.number,
			c.status,
			d.full_name AS driver_full_name,
			d.phone AS driver_phone,
			d.from_city_id AS driver_from_city_id,
			d.to_city_id AS driver_to_city_id,
			d.created_at AS driver_created_at
		FROM
			cars c
		JOIN
			drivers d ON c.driver_id = d.id
		WHERE
			c.id = $1;
	`

	if err := c.db.QueryRow(query, id).Scan(
		&car.ID,
		&car.Model,
		&car.Brand,
		&car.Number,
		&car.Status,
		&car.DriverData.FullName,
		&car.DriverData.Phone,
		&car.DriverData.FromCityID,
		&car.DriverData.ToCityID,
		&car.DriverData.CreatedAt,
	); err != nil {
		fmt.Println("error while scanning user ", err.Error())
		return models.Car{}, err
	}

	return car, nil
}

// TASK 4

func (c carRepo) GetList(req models.GetListRequest) (models.CarsResponse, error) {
	query := `
        SELECT
            cars.id,
            cars.model,
            cars.brand,
            cars.number,
            cars.driver_id,
            cars.status,
            cars.created_at,
            drivers.full_name AS driver_name,
            drivers.phone AS driver_phone,
            drivers.from_city_id AS driver_from_city_id,
            drivers.to_city_id AS driver_to_city_id,
            drivers.created_at AS driver_created_at
        FROM
            cars
        JOIN
            drivers ON cars.driver_id = drivers.id;
    `

	rows, err := c.db.Query(query)
	if err != nil {
		return models.CarsResponse{}, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(
			&car.ID,
			&car.Model,
			&car.Brand,
			&car.Number,
			&car.DriverID,
			&car.Status,
			&car.CreatedAt,
			&car.DriverData.FullName,
			&car.DriverData.Phone,
			&car.DriverData.FromCityID,
			&car.DriverData.ToCityID,
			&car.DriverData.CreatedAt,
		)
		if err != nil {
			return models.CarsResponse{}, fmt.Errorf("error scanning rows: %v", err)
		}

		cars = append(cars, car)
	}

	countQuery := `
        SELECT COUNT(*)
        FROM cars
        JOIN drivers ON cars.driver_id = drivers.id
    `
	var count int
	err = c.db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return models.CarsResponse{}, fmt.Errorf("error executing COUNT query: %v", err)
	}

	return models.CarsResponse{
		Cars:  cars,
		Count: count,
	}, nil
}

func (c carRepo) Update(car models.Car) (string, error) {
	query := `
	UPDATE cars
    SET model = $1, brand = $2, number = $3, driver_id = $4 
    WHERE id = $5;
	`
	if _, err := c.db.Exec(query, car.Model, car.Brand, car.Number, car.DriverID, car.ID); err != nil {
		fmt.Println("error while updating car data ", err.Error())
		return "", err
	}

	return car.ID, nil
}

func (c carRepo) Delete(id string) error {

	query := `delete from cars where id = $1`

	if _, err := c.db.Exec(query, id); err != nil {
		fmt.Println("error while deleting car by id ", err.Error())
		return err
	}
	return nil

}

// TASK 1

func (c carRepo) UpdateCarRoute(updateCarRoute models.UpdateCarRoute) error {

	s, err := c.db.Prepare("UPDATE cars SET departure_time=$1, from_city_id=$2, to_city_id=$3 WHERE car_id=$4")
	if err != nil {
		return err
	}
	defer s.Close()
	updateCarRoute.DepartureTime = time.Now()
	_, err = s.Exec(updateCarRoute.DepartureTime, updateCarRoute.FromCityID, updateCarRoute.ToCityID, updateCarRoute.CarID)
	if err != nil {
		return err
	}

	return nil

}

// TASK 2

func (c carRepo) UpdateCarStatus(updateCarStatus models.UpdateCarStatus) error {

	s, err := c.db.Prepare("UPDATE cars SET status=$1 WHERE id=$2")
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Exec(updateCarStatus.Status, updateCarStatus.ID)
	if err != nil {
		return err
	}
	return nil	

}
