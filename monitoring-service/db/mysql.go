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
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Frequency int    `json:"frequency"`
}

type HistoricalData struct {
	Name         string  `json:"name"`
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

// Update struct
type ResponseData struct {
	Data map[string][]HistoricalData `json:"data"`
	URLs []URL                       `json:"urls"`
}

func (m *Mysql) GetHistoricDataByURLID(urlID int, startDate string, endDate string) (ResponseData, error) {
	responseData := ResponseData{
		Data: make(map[string][]HistoricalData),
	}

	var query string
	var rows *sql.Rows
	var err error

	if urlID != -1 {
		query = `
			SELECT
				u.url as name,
				DATE(h.created_at) AS date,
				AVG(response_time_ms) AS responseTime,
				(SUM(is_alive) / COUNT(*)) * 100 AS uptime
			FROM
				url_health_checks h
				JOIN urls u ON u.id = h.url_id
			WHERE
				url_id = ?
				AND DATE(h.created_at) >= DATE(?)
				AND DATE(h.created_at) <= DATE(?)
			GROUP BY
				DATE(h.created_at)
			ORDER BY
				DATE(h.created_at)`
		rows, err = m.DB.Query(query, urlID, startDate, endDate)
	} else {
		query = `
		SELECT
			u.url as name,
			DATE(h.created_at) AS date,
			AVG(h.response_time_ms) AS responseTime,
			(SUM(h.is_alive) / COUNT(*)) * 100 AS uptime
		FROM
			urls u
			JOIN url_health_checks h ON u.id = h.url_id
		WHERE
			DATE(h.created_at) >= DATE(?)
			AND DATE(h.created_at) <= DATE(?)
		GROUP BY
			u.url,
			DATE(h.created_at)
		ORDER BY
			u.url,
			DATE(h.created_at)`
		rows, err = m.DB.Query(query, startDate, endDate)
	}

	fmt.Println(query)

	if err != nil {
		return responseData, err
	}

	defer rows.Close()

	// Process rows into grouped data
	urlNames := make(map[string]bool)
	for rows.Next() {
		var d HistoricalData
		err := rows.Scan(&d.Name, &d.Date, &d.ResponseTime, &d.Uptime)
		if err != nil {
			return responseData, err
		}

		responseData.Data[d.Name] = append(responseData.Data[d.Name], d)
		urlNames[d.Name] = true
	}

	return responseData, nil
}

func (m *Mysql) Close() {
	m.DB.Close()
}
