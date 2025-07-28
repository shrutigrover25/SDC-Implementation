package db

import (
	"log"

	"gorm.io/driver/postgres"

	jobs "mercor/internal/domain/jobs"
	paymentLineItem "mercor/internal/domain/paymentLineItem"
	timelog "mercor/internal/domain/timelog"

	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := "host=localhost user=postgres password=Shruti@25 dbname=mercortwo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(
		&jobs.Job{},
		&timelog.Timelog{},
		&paymentLineItem.PaymentLineItem{},
	)
	if err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	return db
}
