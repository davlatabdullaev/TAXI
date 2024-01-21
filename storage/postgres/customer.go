package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type customerRepo struct {
	db *sql.DB
}

func NewCustomerRepo(db *sql.DB) storage.ICustomerRepo {
	return customerRepo{
		db,
	}
}

func (c customerRepo) Create(customer models.CreateCustomer) (string, error) {

	uid := uuid.New()

	if _, err := c.db.Exec(`
	 insert into customers values ($1, $2, $3, $4)
	 `,
		uid,
		customer.FullName,
		customer.Phone,
		customer.Email,
	); err != nil {
		fmt.Println("error while inserting data", err.Error())
		return "", err
	}

	return uid.String(), nil

}

func (c customerRepo) Get(id string) (models.Customer, error) {

	customer := models.Customer{}

	query := `
		select id, full_name, phone, email, created_at from customers where id = $1
`
	if err := c.db.QueryRow(query, id).Scan(
		&customer.ID,
		&customer.FullName,
		&customer.Phone,
		&customer.Email,
		&customer.CreatedAt,
	); err != nil {
		fmt.Println("error while scanning user", err.Error())
		return models.Customer{}, err
	}

	return customer, nil

}

func (c customerRepo) GetList(req models.GetListRequest) (models.CustomersResponse, error) {

	var (
		customers         = []models.Customer{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
	)

	countQuery = `
	SELECT count(1) from customers `

	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of users", err.Error())
		return models.CustomersResponse{}, err
	}

	query = `
	SELECT id, full_name, phone, email, created_at
		FROM customers
			`

	query += ` LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.CustomersResponse{}, err
	}

	for rows.Next() {
		customer := models.Customer{}

		if err = rows.Scan(
			&customer.ID,
			&customer.FullName,
			&customer.Phone,
			&customer.Email,
			&customer.CreatedAt,
		); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.CustomersResponse{}, err
		}

		customers = append(customers, customer)
	}

	return models.CustomersResponse{
		Customers: customers,
		Count:     count,
	}, nil

}

func (c customerRepo) Update(customer models.Customer) (string, error) {
	query := `
	update customers 
		set full_name = $1, phone = $2, email = $3
			where id = $4`

	if _, err := c.db.Exec(query, customer.FullName, customer.Phone, customer.Email, customer.ID); err != nil {
		fmt.Println("error while updating customer data", err.Error())
		return "", err
	}

	return customer.ID, nil

}

func (c customerRepo) Delete(id string) error {
	query := `
	delete from customers
		where id = $1
`
	if _, err := c.db.Exec(query, id); err != nil {
		fmt.Println("error while deleting customer by id", err.Error())
		return err
	}

	return nil
}
