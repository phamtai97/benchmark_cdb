package storage

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

//PostgreSQL struct
type PostgreSQL struct {
	db       *sql.DB
	user     string
	password string
	host     string
	port     int
	dbName   string
	sslmode  bool
}

//ConnectDB connect db
func (PostgreSQL *PostgreSQL) ConnectDB() error {
	var mode string
	if PostgreSQL.sslmode {
		mode = "disable"
	}

	dataSourceName := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		PostgreSQL.user, PostgreSQL.password, PostgreSQL.host, PostgreSQL.port, PostgreSQL.dbName, mode)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return errors.New(err.Error())
	}

	db.SetMaxOpenConns(MAX_OPEN_CONNS)
	db.SetMaxIdleConns(MAX_IDLE_CONNS)

	PostgreSQL.db = db
	log.Printf("[Storage] Connect DB success host: %s port: %d", PostgreSQL.host, PostgreSQL.port)
	return nil
}

//CloseDB close the database
func (PostgreSQL *PostgreSQL) CloseDB() error {
	log.Println("[Storage] Close database PostgreSQL")
	return PostgreSQL.db.Close()
}

//CreateTable create table
func (PostgreSQL *PostgreSQL) CreateTable() error {
	if _, err := PostgreSQL.db.Exec(
		"CREATE TABLE IF NOT EXISTS accounts" +
			"(id varchar(18) PRIMARY KEY, balance bigint)"); err != nil {
		return err
	}
	return nil
}

//InitStoragePostgre init storage
func InitStoragePostgre(user string, password string, host string, port int, dbName string, sslmode bool) (*PostgreSQL, error) {
	PostgreSQL := PostgreSQL{
		db:       nil,
		user:     user,
		password: password,
		host:     host,
		port:     port,
		dbName:   dbName,
		sslmode:  sslmode,
	}

	if err := PostgreSQL.ConnectDB(); err != nil {
		return &PostgreSQL, err
	}

	if err := PostgreSQL.CreateTable(); err != nil {
		return &PostgreSQL, err
	}

	return &PostgreSQL, nil
}

//GetDB get db
func (PostgreSQL *PostgreSQL) GetDB() *sql.DB {
	return PostgreSQL.db
}
