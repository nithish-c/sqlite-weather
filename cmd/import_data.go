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
	Db_file           = filepath.Join(".", "weather", "weather.db")
	Measured_data_10k = filepath.Join(".", "weather", "measurements_10k.txt")
	// Measured_data_1b   = filepath.Join(".", "weather", "measurements_1b.txt")
	// Measured_data_100m = filepath.Join(".", "weather", "measurements_100m.txt")
	// Measured_data_10m = filepath.Join(".", "weather", "measurements_10m.txt")
	Measured_data_5m = filepath.Join(".", "weather", "measurements_5m.txt")
)

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func OpenDB(db_name string) (*sql.DB, error) {
	os.Remove(db_name)
	_, err := os.Stat(db_name)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("File does not exist.")
			_, e := os.Create(db_name)
			errorCheck(e)
		}
	}
	return sql.Open("sqlite3", db_name)
}

func NewWeather(city string, temp float64) w.Weather {
	return w.Weather{
		City: city,
		Temp: temp,
	}
}

func ImportData(db *w.RepositorySQLite, file string) {
	// Open file
	f, err := os.Open(file)
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
