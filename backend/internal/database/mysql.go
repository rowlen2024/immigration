package database

import (
	"crypto/tls"
	"log"
	"time"

	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/model"

	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL(cfg *config.Config) (*gorm.DB, error) {
	mysql.RegisterTLSConfig("skip-verify", &tls.Config{InsecureSkipVerify: true})

	db, err := gorm.Open(gormmysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("MySQL connected successfully")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Requirement{},
		&model.CostItem{},
		&model.TimelinePhase{},
		&model.Milestone{},
		&model.ProjectAdvantage{},
		&model.FAQ{},
		&model.Case{},
		&model.Lawyer{},
		&model.Page{},
		&model.Lead{},
		&model.Media{},
		&model.Navigation{},
		&model.HomeConfig{},
		&model.OperationLog{},
	)
}
