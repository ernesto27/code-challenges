package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	DB *sql.DB
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

func (m *Mysql) UpdateURLFrequency(url string, frequency int) error {
	_, err := m.DB.Exec("UPDATE urls SET frequency = ? WHERE url = ?", frequency, url)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mysql) AddURLHealthCheck(urlID int, statusCode int, responseTime int, isAlive int) error {
	_, err := m.DB.Exec("INSERT INTO url_health_checks (url_id, status_code, response_time_ms, is_alive) VALUES (?, ?, ?, ?)", urlID, statusCode, responseTime, isAlive)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mysql) Close() {
	m.DB.Close()
}
