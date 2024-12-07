package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	DB *sql.DB
}

type URL struct {
	ID        int
	URL       string
	Frequency int
}

type HistoricalData struct {
	Date         string  `json:"date"`
	ResponseTime float64 `json:"responseTime"`
	Uptime       float64 `json:"uptime"`
}

func NewMysql(host, user, password, port, database string) (*Mysql, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.New("error connecting to the database")
	}

	return &Mysql{
		DB: db,
	}, nil
}

func (m *Mysql) CreateURL(url string, frequency int) error {
	_, err := m.DB.Exec("INSERT INTO urls (url, frequency) VALUES (?, ?)", url, frequency)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mysql) GetURLs() ([]URL, error) {
	rows, err := m.DB.Query("SELECT id, url, frequency FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		err := rows.Scan(&url.ID, &url.URL, &url.Frequency)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	return urls, nil
}

func (m *Mysql) UpdateURLFrequency(url string, frequency int) error {
	_, err := m.DB.Exec("UPDATE urls SET frequency = ? WHERE url = ?", frequency, url)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mysql) CreateURLHealthCheck(urlID int, statusCode int, responseTime int, isAlive int) error {
	_, err := m.DB.Exec("INSERT INTO url_health_checks (url_id, status_code, response_time_ms, is_alive) VALUES (?, ?, ?, ?)", urlID, statusCode, responseTime, isAlive)
	if err != nil {
		return err
	}

	return nil
}

type ResponseData struct {
	Name string           `json:"name"`
	Data []HistoricalData `json:"data"`
}

func (m *Mysql) GetHistoricDataByURLID(urlID int) (ResponseData, error) {
	responseData := ResponseData{}

	row := m.DB.QueryRow("SELECT url FROM urls WHERE id = ?", urlID)
	var url string

	err := row.Scan(&url)
	if err != nil {
		return responseData, err
	}

	rows, err := m.DB.Query(
		`SELECT 
			DATE(created_at) AS date,
			AVG(response_time_ms) AS responseTime,
			(SUM(is_alive) / COUNT(*)) * 100 AS uptime
		FROM 
			url_health_checks
		WHERE 
			url_id = ?
		GROUP BY 
			DATE(created_at)
		ORDER BY 
			DATE(created_at) 
	`, urlID)

	if err != nil {
		return responseData, err
	}

	defer rows.Close()

	var data []HistoricalData
	for rows.Next() {
		var d HistoricalData
		err := rows.Scan(&d.Date, &d.ResponseTime, &d.Uptime)
		if err != nil {
			return responseData, err
		}

		data = append(data, d)
	}

	responseData.Name = url
	responseData.Data = data

	return responseData, nil

}

func (m *Mysql) Close() {
	m.DB.Close()
}
