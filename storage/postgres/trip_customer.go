package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type tripCustomerRepo struct {
	db *sql.DB
}

func NewTripCustomerRepo(db *sql.DB) storage.ITripCustomerRepo {
	return &tripCustomerRepo{
		db: db,
	}
}

func (c tripCustomerRepo) Create(req models.CreateTripCustomer) (string, error) {

	uid := uuid.New()

	if _, err := c.db.Exec(`insert into 
			trip_customers 
			(id, trip_id, customer_id) values($1, $2, $3)
			`,
		uid,
		req.TripID,
		req.CustomerID,
	); err != nil {
		fmt.Println("error while inserting data", err.Error())
		return "", err
	}

	return uid.String(), nil

}

func (c tripCustomerRepo) Get(id string) (models.TripCustomer, error) {
	trip := models.TripCustomer{}
	query := `SELECT tr.id, tr.trip_id, tr.customer_id, 
       				 c.id as customer_id,c.full_name as customer_name, c.phone as customer_phone, 
       				 c.email as customer_email, c.created_at as customer_date,
       				 tr.created_at
					FROM trip_customers as tr 
					LEFT JOIN customers as c ON tr.customer_id = c.id WHERE tr.id = $1`

	if err := c.db.QueryRow(query, id).Scan(
		&trip.ID,
		&trip.TripID,
		&trip.CustomerID,
		&trip.CustomerData.ID,
		&trip.CustomerData.FullName,
		&trip.CustomerData.Phone,
		&trip.CustomerData.Email,
		&trip.CustomerData.CreatedAt,
		&trip.CreatedAt,
	); err != nil {
		fmt.Println("error is while scanning trip customer", err.Error())
		return models.TripCustomer{}, err
	}
	return trip, nil
}

func (c tripCustomerRepo) GetList(req models.GetListRequest) (models.TripCustomersResponse, error) {
	var (
		page              = req.Page
		offset            = (page - 1) * req.Limit
		tripCustomers     = []models.TripCustomer{}
		query, countQuery string
		count             = 0
	)

	countQuery = `SELECT count(1) FROM trip_customers`
	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning count", err.Error())
		return models.TripCustomersResponse{}, err
	}

	query = `SELECT tr.id, tr.trip_id, tr.customer_id, 
       				 c.id as customer_id,c.full_name as customer_name, c.phone as customer_phone, 
       				 c.email as customer_email, c.created_at as customer_date,
       				 tr.created_at
					FROM trip_customers as tr 
					LEFT JOIN customers as c ON tr.customer_id = c.id `
	query += ` LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting trip customers", err.Error())
		return models.TripCustomersResponse{}, err
	}

	for rows.Next() {
		trip := models.TripCustomer{}
		if err = rows.Scan(
			&trip.ID, &trip.TripID, &trip.CustomerID,
			&trip.CustomerData.ID, &trip.CustomerData.FullName, &trip.CustomerData.Phone, &trip.CustomerData.Email,
			&trip.CustomerData.CreatedAt, &trip.CreatedAt,
		); err != nil {
			fmt.Println("error is while scanning rows", err.Error())
			return models.TripCustomersResponse{}, err
		}
		tripCustomers = append(tripCustomers, trip)
	}

	return models.TripCustomersResponse{
		TripCustomers: tripCustomers,
		Count:         count,
	}, nil
}

func (c tripCustomerRepo) Update(req models.TripCustomer) (string, error) {
	query := `UPDATE trip_customers SET customer_id = $1 WHERE id = $2`
	if _, err := c.db.Exec(query, req.CustomerID, req.ID); err != nil {
		fmt.Println("error is while updating trip customer", err.Error())
		return "", err
	}
	return req.ID, nil
}

func (c *tripCustomerRepo) Delete(id string) error {
	query := `DELETE FROM trip_customers WHERE id = $1`

	if _, err := c.db.Exec(query, id); err != nil {
		fmt.Println("error is while deleting trip customer", err.Error())
		return err
	}
	return nil
}
