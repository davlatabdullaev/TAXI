package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type cityRepo struct {
	db *sql.DB
}

func NewCityRepo(db *sql.DB) storage.ICityRepo {
	return cityRepo{
		db,
	}
}

func (c cityRepo) Create(city models.CreateCity) (string, error) {
	uid := uuid.New()

	query := `
     insert into cities values ($1, $2) 
   `

	if _, err := c.db.Exec(query, uid, city.Name); err != nil {
		fmt.Println("error while inserting data", err.Error())
		return "", err
	}

	return uid.String(), nil

}

func (c cityRepo) Get(id string) (models.City, error) {

	city := models.City{}

	query := `select id, name, created_at from cities where id = $1`

	if err := c.db.QueryRow(query, id).Scan(
		&city.ID,
		&city.Name,
		&city.CreatedAt,
	); err != nil {
		fmt.Println("error while scanning city", err.Error())
		return models.City{}, err
	}

	return city, nil

}

func (c cityRepo) GetList(req models.GetListRequest) (models.CitiesResponse, error) {

	var (
		cities            = []models.City{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
	)

	countQuery = `select count(1) from cities`

	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of cities", err.Error())
		return models.CitiesResponse{}, err
	}

	query = `
	select id, 
	name,
	 created_at from cities
	`

	query += ` limit $1 offset $2`

	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.CitiesResponse{}, err
	}

	for rows.Next() {
		city := models.City{}

		if err = rows.Scan(
			&city.ID,
			&city.Name,
			&city.CreatedAt,
		); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.CitiesResponse{}, err
		}

		cities = append(cities, city)

	}

	return models.CitiesResponse{
		Cities: cities,
		Count:  count,
	}, nil
}

func (c cityRepo) Update(city models.City) (string, error) {

	query := `update
	cities 
	set name = $1 where 
	id = $2`

	if _, err := c.db.Exec(query, city.Name, city.ID); err != nil {
		fmt.Println("error while updating city data ", err.Error())
		return "", err
	}

	return city.ID, nil

}

func (c cityRepo) Delete(id string) error {

	query := `
     delete from cities 
     where id = $1 
  `
	if _, err := c.db.Exec(query, id); err != nil {
		fmt.Println("error while deleeting city by id", err.Error())
		return err
	}

	return nil

}
