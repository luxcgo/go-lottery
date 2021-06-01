package main

import (
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "secret"
	dbname   = "luxcgo_gallery"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;uniqueIndex"`
}

func main() {}
