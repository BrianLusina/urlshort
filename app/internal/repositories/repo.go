package repositories

import (
	"fmt"
	"log"
	"quote/api/app/config"
	"quote/api/app/internal/repositories/models"
	"quote/api/app/internal/repositories/quotesrepo"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db         *gorm.DB
	quotesRepo *quotesrepo.QuotesRepo
}

func NewRepository(config config.DatabaseConfig) *repository {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Host, config.User, config.Password, config.Database, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("DB Connection failed with err: %v", err)
	}

	if err = db.AutoMigrate(&models.Quote{}); err != nil {
		log.Fatalf("AutoMigration failed with err: %v", err)
	}

	return &repository{
		db:         db,
		quotesRepo: quotesrepo.NewQuotesRepo(db),
	}
}

func (r repository) GetQuotesRepo() *quotesrepo.QuotesRepo {
	return r.quotesRepo
}
