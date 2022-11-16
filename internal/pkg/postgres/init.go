package postgres

import (
	"fmt"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(config Postgres) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open db")
	}

	return db, nil
}
