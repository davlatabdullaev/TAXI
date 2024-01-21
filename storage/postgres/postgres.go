package postgres

import (
	"city2city/config"
	"city2city/storage"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func New(cfg config.Config) (storage.IStorage, error) {
	url := fmt.Sprintf(`host = %s port = %s user = %s password = %s database = %s sslmode=disable`,
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return Store{}, err
	}

	return Store{
		db: db,
	}, nil
}

func (s Store) CloseDB() {
	s.db.Close()
}

func (s Store) City() storage.ICityRepo {
	return NewCityRepo(s.db)
}

func (s Store) Customer() storage.ICustomerRepo {
	return NewCustomerRepo(s.db)
}

func (s Store) Driver() storage.IDriverRepo {
	return NewDriverRepo(s.db)
}

func (s Store) Car() storage.ICarRepo {
	return NewCarRepo(s.db)
}

func (s Store) Trip() storage.ITripRepo {
	return NewTripRepo(s.db)
}
func (s Store) TripCustomer() storage.ITripCustomerRepo {
	return NewTripCustomerRepo(s.db)
}