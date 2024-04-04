package main

import (
	"log"
	c "simple-db/cmd"
	w "simple-db/weather"

	_ "github.com/mattn/go-sqlite3"
)

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	db, err := c.OpenDB(c.Db_file)
	errorCheck(err)
	defer db.Close()

	weatherRepository := w.NewRepositorySQLite(db)
	errCreate := weatherRepository.Migrate()
	errorCheck(errCreate)

	c.ImportData(weatherRepository)
	c.UniqueCitiesDB(weatherRepository)
}
