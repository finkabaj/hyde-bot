package db

import ()

type Database interface {
	Connect(credentials *DatabaseCredentials) error
	Close()
	Status() error
}

type DatabaseCredentials struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}
