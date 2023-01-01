package ipdb

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/oschwald/geoip2-golang"
)

type DBResponse struct {
	City_EN      string `json:"city_en"`
	Country_EN   string `json:"country_en"`
	Country_DE   string `json:"country_de"`
	Country_Code string `json:"country_code"`
	Continent_EN string `json:"continent_en"`
	Continent_DE string `json:"continent_de"`
	Longitude    string `json:"longitude"`
	Latitude     string `json:"latitude"`
	Subdivision  string `json:"subdivision"`
	IP           string `json:"ip"`
}

type IPDB struct {
	db *geoip2.Reader
}

func downloadCurrentIPDB() (string, error) {
	// get current month and year
	monthYear := getCurrentMonthYear()

	currentFilename := "dbip-country-lite-" + monthYear + ".mmdb"
	if _, err := os.Stat(currentFilename); err == nil {
		fmt.Println("File already exists, skipping download")
		return currentFilename, err
	}
	url := "https://download.db-ip.com/free/dbip-city-lite-" + monthYear + ".mmdb.gz"

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create a new file to save the downloaded data
	file, err := os.Create("dbip-country-lite-" + monthYear + ".mmdb.gz")
	if err != nil {
		return "", err
	}
	defer file.Close()
	fmt.Println("Downloading file...")

	// Copy the data from the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	db, err := unzip("dbip-country-lite-" + monthYear + ".mmdb.gz")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Download complete! File saved to", db)

	return db, nil

}

func NewDB() (*IPDB, error) {

	db_filename, err := downloadCurrentIPDB()
	if err != nil {
		return nil, err
	}

	db, err := geoip2.Open(db_filename)
	if err != nil {
		log.Fatal(err)
	}
	return &IPDB{
		db: db,
	}, nil
}

func unzip(filename string) (string, error) {
	// Open the gzip file
	gzipFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer gzipFile.Close()

	// Create a new gzip reader
	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	// Create a new file to write the decompressed data
	file, err := os.Create(filename[:len(filename)-3]) // remove ".gz
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Copy the decompressed data from the gzip reader to the file
	_, err = io.Copy(file, gzipReader)
	if err != nil {
		return "", err
	}
	return filename[:len(filename)-3], nil

}

func getCurrentMonthYear() string {
	now := time.Now()

	// Format the month and year as "YYYY-MM"
	monthYear := now.Format("2006-01")
	return monthYear

}

func (db *IPDB) LookUpIP(input string) (DBResponse, error) {
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(input)
	record, err := db.db.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	response := DBResponse{
		City_EN:      record.City.Names["en"],
		Country_Code: record.Country.IsoCode,
		Country_DE:   record.Country.Names["de"],
		Country_EN:   record.Country.Names["en"],
		Continent_EN: record.Continent.Names["en"],
		Continent_DE: record.Continent.Names["de"],
		Longitude:    fmt.Sprintf("%f", record.Location.Longitude),
		Latitude:     fmt.Sprintf("%f", record.Location.Latitude),
		IP:           input,
	}
	if len(record.Subdivisions) > 0 {
		response.Subdivision = record.Subdivisions[0].Names["en"]
	}

	return response, nil
}
func (db *IPDB) Close() {
	db.db.Close()
}
