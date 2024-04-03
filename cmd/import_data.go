package cmd

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	w "simple-db/weather"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db_file       = filepath.Join(".", "weather", "weather.db")
	measured_data = filepath.Join(".", "weather", "measurements.txt")
)

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func OpenDB() (*sql.DB, error) {
	// os.Remove(db_file)
	_, err := os.Stat(db_file)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("File does not exist.")
			_, e := os.Create(db_file)
			errorCheck(e)
		}
	}
	return sql.Open("sqlite3", db_file)
}

func NewWeather(city string, temp float64) w.Weather {
	return w.Weather{
		City: city,
		Temp: temp,
	}
}

func ImportData(db *w.RepositorySQLite) {
	// Open file
	f, err := os.Open(measured_data)
	errorCheck(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, ";")
		temp, e := strconv.ParseFloat(data[1], 64)
		errorCheck(e)
		w := NewWeather(data[0], temp)
		_, err = db.Create(w)
		errorCheck(err)
	}
}
