package main

import (
	"database/sql"
	"log"
	"os"
	w "simple-db/weather"

	_ "github.com/mattn/go-sqlite3"
)

var weather1 = w.Weather{
	Id:   1,
	City: "Jakarta",
	Temp: 30.0,
}

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	// Check if file exists, create new if doesn't exist. To be done in later.
	os.Remove("./weather.db")
	return sql.Open("sqlite3", "./weather.db")
}
func main() {

	db, err := openDB()
	errorCheck(err)
	defer db.Close()

	weatherRepository := w.NewRepositorySQLite(db)
	errCreate := weatherRepository.Migrate()
	errorCheck(errCreate)
	_, err_Update := weatherRepository.Create(weather1)
	errorCheck(err_Update)
	_, err_Update = weatherRepository.Create(newWeather("Bandung", 25.0))
	errorCheck(err_Update)
}

func newWeather(city string, temp float64) w.Weather {
	return w.Weather{
		City: city,
		Temp: temp,
	}
}
