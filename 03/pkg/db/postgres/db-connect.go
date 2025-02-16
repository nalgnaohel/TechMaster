package postgres

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"03/config"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"gorm.io/driver/postgres"
)

// NewPsqlDB Return new Postgresql db instance
func NewPsqlDB(c *config.Config) (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.DbName,
		c.Postgres.Password,
	)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(c.Postgres.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.Postgres.ConnMaxLifetime) * time.Second)
	sqlDB.SetMaxIdleConns(c.Postgres.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(time.Duration(c.Postgres.ConnMaxIdleTime) * time.Second)
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
