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

func (c *tripCustomerRepo) Create(req models.CreateTripCustomer) (string, error) {

	uid := uuid.New()

	if _, err := c.db.Exec(`insert into 
			trip_customers values ($1, $2, $3)
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

// TASK 7

func (c *tripCustomerRepo) Get(id string) (models.TripCustomer, error) {

	tripCustomer := models.TripCustomer{}

	query := `
		select id, trip_id, customer_id, customer_data, created_at from trip_customers where id = $1 
`
	if err := c.db.QueryRow(query, id).Scan(
		&tripCustomer.ID,
		&tripCustomer.TripID,
		&tripCustomer.CustomerID,
		&tripCustomer.CustomerData,
		&tripCustomer.CreatedAt,
	); err != nil {
		fmt.Println("error while scanning trip customer", err.Error())
		return models.TripCustomer{}, err
	}

	return models.TripCustomer{}, nil

}

// TASK 8

func (c *tripCustomerRepo) GetList(req models.GetListRequest) (models.TripCustomersResponse, error) {
	var (
		tripCustomers     = []models.TripCustomer{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
	)

	countQuery = `
		SELECT count(1) from trip_customers `

	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of users", err.Error())
		return models.TripCustomersResponse{}, err
	}

	query = `
		SELECT id, trip_id, customer_id, customer_data, created_at
			FROM trip_customers
			    `

	query += ` LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.TripCustomersResponse{}, err
	}

	for rows.Next() {
		tripCustomer := models.TripCustomer{}

		if err = rows.Scan(
			&tripCustomer.ID,
			&tripCustomer.TripID,
			&tripCustomer.CustomerID,
			&tripCustomer.CustomerData,
			&tripCustomer.CreatedAt,
		); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.TripCustomersResponse{}, err
		}

		tripCustomers = append(tripCustomers, tripCustomer)
	}

	return models.TripCustomersResponse{
		TripCustomers: tripCustomers,
		Count:         count,
	}, nil

}

func (c *tripCustomerRepo) Update(req models.TripCustomer) (string, error) {
	query := `
	update trip_customers 
		set trip_id = $1, customer_id = $2, customer_data = $3
			where id = $4`

	if _, err := c.db.Exec(query, req.TripID, req.CustomerID, req.CustomerData, req.ID); err != nil {
		fmt.Println("error while updating trip customer data", err.Error())
		return "", err
	}
	return req.ID, nil

}

func (c *tripCustomerRepo) Delete(id string) error {
	query := `
	delete from trip_customers
		where id = $1
`
	if _, err := c.db.Exec(query, id); err != nil {

		fmt.Println("error while deleting trip customer by id", err.Error())

		return err
	}

	return nil
}
