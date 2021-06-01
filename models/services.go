package models

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Services struct {
	Gallery GalleryService
	db      *gorm.DB
}

// ServicesConfig is really just a function, but I find using
// types like this are easier to read in my source code.
type ServicesConfig func(*Services) error

// NewServices now will accept a list of config functions to
// run. Each function will accept a pointer to the current
// Services object as its only argument and will edit that
// object inline and return an error if there is one. Once
// we have run all configs we will return the Services object.
func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	// For each ServicesConfig function...
	for _, cfg := range cfgs {
		// Run the function passing in a pointer to our Services
		// object and catching any errors
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	// Then finally return the result
	return &s, nil
}

func WithGorm(connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(postgres.Open(connectionInfo))
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithNewLogger(mode bool) ServicesConfig {
	return func(s *Services) error {
		if mode {
			newLogger := logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold: time.Second, // Slow SQL threshold
					LogLevel:      logger.Info, // Log level
					Colorful:      true,        // Disable color
				},
			)
			s.db.Logger = newLogger
		}
		return nil
	}
}

func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

// AutoMigrate will attempt to automatically migrate all tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&Gallery{})
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.Migrator().DropTable(&Gallery{})
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
