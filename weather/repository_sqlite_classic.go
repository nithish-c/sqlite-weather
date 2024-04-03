package weather

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

type RepositorySQLite struct {
	db *sql.DB
}

func NewRepositorySQLite(db *sql.DB) *RepositorySQLite {
	return &RepositorySQLite{db: db}
}

func (r *RepositorySQLite) Migrate() error {
	query := `CREATE TABLE IF NOT EXISTS weather (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		city TEXT NOT NULL,
		temp REAL NOT NULL
	);`
	_, err := r.db.Exec(query)
	return err
}

func (r *RepositorySQLite) Create(weather Weather) (*Weather, error) {
	query := `INSERT INTO weather (city, temp) VALUES (?, ?);`
	result, err := r.db.Exec(query, weather.City, weather.Temp)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqlite3.ErrConstraintUnique, sqliteErr.ExtendedCode) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	weather.Id = id
	return &weather, nil
}

func (r *RepositorySQLite) All() ([]Weather, error) {
	query := `SELECT id, city, temp FROM weather;`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weathers []Weather
	for rows.Next() {
		var weather Weather
		err := rows.Scan(&weather.Id, &weather.City, &weather.Temp)
		if err != nil {
			return nil, err
		}
		weathers = append(weathers, weather)
	}
	return weathers, nil
}

func (r *RepositorySQLite) GetByID(id int64) (*Weather, error) {
	query := `SELECT * FROM weather WHERE id = ?;`
	row := r.db.QueryRow(query, id)

	var data_point Weather
	err := row.Scan(&data_point.Id, &data_point.City, &data_point.Temp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &data_point, nil
}

func (r *RepositorySQLite) GetByCity(city string) ([]Weather, error) {
	query := `SELECT * FROM weather WHERE city = ?;`
	rows, err := r.db.Query(query, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []Weather
	for rows.Next() {
		var weather Weather
		err := rows.Scan(&weather.Id, &weather.City, &weather.Temp)
		if err != nil {
			return nil, err
		}
		cities = append(cities, weather)
	}
	return cities, nil
}

func (r *RepositorySQLite) Update(weather Weather) (*Weather, error) {
	query := `UPDATE weather SET city = ?, temp = ? WHERE id = ?;`
	result, err := r.db.Exec(query, weather.City, weather.Temp, weather.Id)
	if err != nil {
		return nil, ErrUpdate
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, ErrNotFound
	}
	return &weather, nil
}

func (r *RepositorySQLite) Delete(id int64) error {
	query := `DELETE FROM weather WHERE id = ?;`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return ErrDelete
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *RepositorySQLite) DropTable() error {
	query := `DROP TABLE IF EXISTS weather;`
	_, err := r.db.Exec(query)
	return err
}
