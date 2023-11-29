package internal

import "fmt"

type DatabaseConnectionStr struct {
	User     string
	Host     string
	Port     string
	Database string
	Password string
}

func (conn DatabaseConnectionStr) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", conn.Host, conn.Port, conn.User, conn.Database, conn.Password)
}
