package cmd

import (
	"path/filepath"
	w "simple-db/weather"
)

var Unique_cites = filepath.Join(".", "weather", "unique.db")

func UniqueCitiesDB(main *w.RepositorySQLite) {
	db, err1 := OpenDB(Unique_cites)
	errorCheck(err1)
	defer db.Close()

	// r := w.NewRepositorySQLite(db)
	cities, err2 := w.UniqueCities(main)
	errorCheck(err2)

	_, err3 := db.Exec("CREATE TABLE IF NOT EXISTS unique_cities (city TEXT NOT NULL);")
	errorCheck(err3)
	for _, city := range cities {
		_, err4 := db.Exec("INSERT INTO unique_cities (city) VALUES (?)", city)
		errorCheck(err4)
	}
}
