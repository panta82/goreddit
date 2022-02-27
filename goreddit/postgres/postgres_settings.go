package postgres

import "fmt"

type Settings struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func NewSettings() *Settings {
	return &Settings{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		Name:     "goreddit",
	}
}

func (ss *Settings) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		ss.User, ss.Password, ss.Host, ss.Port, ss.Name)
}
