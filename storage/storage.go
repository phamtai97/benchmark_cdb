package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	MAX_OPEN_CONNS = 200
	MAX_IDLE_CONNS = 50
)

//Storage storage
type Storage interface {
	ConnectDB() error
	CloseDB() error
	CreateTable() error
	GetDB() *sql.DB
}

//CockroachDB struct
type CockroachDB struct {
	db      *sql.DB
	user    string
	host    string
	port    int
	dbName  string
	sslmode bool
}

//ConnectDB connect db
func (cockroachDB *CockroachDB) ConnectDB() error {
	var mode string
	if cockroachDB.sslmode {
		mode = "disable"
	}

	dataSourceName := fmt.Sprintf("user=%s host=%s port=%d dbname=%s sslmode=%s",
		cockroachDB.user, cockroachDB.host, cockroachDB.port, cockroachDB.dbName, mode)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return errors.New(err.Error())
	}

	db.SetMaxOpenConns(MAX_OPEN_CONNS)
	db.SetMaxIdleConns(MAX_IDLE_CONNS)

	cockroachDB.db = db
	log.Printf("[Storage] Connect DB host: %s port: %d", cockroachDB.host, cockroachDB.port)
	return nil
}

//CloseDB close the database
func (cockroachDB *CockroachDB) CloseDB() error {
	log.Println("[Storage] Close database cockroachDB")
	return cockroachDB.db.Close()
}

//CreateTable create table
func (cockroachDB *CockroachDB) CreateTable() error {
	if _, err := cockroachDB.db.Exec(
		"CREATE TABLE IF NOT EXISTS benchmark.accounts" +
			"(id varchar(18) PRIMARY KEY, balance bigint)"); err != nil {
		return err
	}
	return nil
}

//InitStorage init storage
func InitStorage(user string, host string, port int, dbName string, sslmode bool) (*CockroachDB, error) {
	cockroachDB := CockroachDB{
		db:      nil,
		user:    user,
		host:    host,
		port:    port,
		dbName:  dbName,
		sslmode: sslmode,
	}

	if err := cockroachDB.ConnectDB(); err != nil {
		return &cockroachDB, err
	}

	if err := cockroachDB.CreateTable(); err != nil {
		return &cockroachDB, err
	}

	return &cockroachDB, nil
}

//GetDB get db
func (cockroachDB *CockroachDB) GetDB() *sql.DB {
	return cockroachDB.db
}

