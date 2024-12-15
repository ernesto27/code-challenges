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
	ID                  int    `json:"id"`
	URL                 string `json:"url"`
	Frequency           int    `json:"frequency"`
	CurrentFrequency    int    `json:"currentFrequency"`
	AttempsFails        int    `json:"attempsFails"`
	CurrentAttempsFails int    `json:"currentAttempsFail"`
}

type HistoricalData struct {
	Name         string  `json:"name"`
	Date         string  `json:"date"`
	ResponseTime float64 `json:"responseTime"`
	Uptime       float64 `json:"uptime"`
}

type Notification struct {
	ID      int
	Message string
	Sent    bool
	Email   string
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
	rows, err := m.DB.Query("SELECT id, url, frequency, attemps_fails FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		err := rows.Scan(&url.ID, &url.URL, &url.Frequency, &url.AttempsFails)
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

func (m *Mysql) CreateURLHealthCheck(urlID int, statusCode int, responseTimeHead int, responseTimeTTFB int, responseTimeGet int, isAlive int) error {
	_, err := m.DB.Exec(`
		INSERT INTO url_health_checks 
			(url_id, status_code, response_time_ms_head, response_time_ms_ttfb, response_time_ms_get, is_alive) 
			VALUES (?, ?, ?, ?, ?, ?)`,
		urlID, statusCode, responseTimeHead, responseTimeTTFB, responseTimeGet, isAlive)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mysql) CreateNotifiction(urlID int, message string) error {
	_, err := m.DB.Exec("INSERT INTO notifications (url_id, message) VALUES (?, ?)", urlID, message)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mysql) GetNotifications() ([]Notification, error) {
	rows, err := m.DB.Query(`
		SELECT n.id, n.message, n.sent, users.email
		FROM notifications n
		JOIN urls ON n.url_id = urls.id
		JOIN users ON urls.user_id = users.id
		WHERE n.sent = 0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []Notification{}
	for rows.Next() {
		var notification Notification
		err := rows.Scan(&notification.ID, &notification.Message, &notification.Sent, &notification.Email)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (m *Mysql) UpdateNotificationSent(id int) error {
	_, err := m.DB.Exec("UPDATE notifications SET sent = 1 WHERE id = ?", id)
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
