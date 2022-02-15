package database

import (
	"log"

	"github.com/NAziz21/newCreditProduct/pkg/models"
	"github.com/NAziz21/newCreditProduct/settings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error

	DB, err = gorm.Open(settings.Config.Driver, settings.Config.DB)
	if err != nil {
		panic("Не удалось подключиться к БД")
	}

	log.Println("DB connected!")

	postgresDB := DB.DB()
	postgresDB.SetMaxOpenConns(100)
	DB.AutoMigrate(&models.Clients{}, &models.Debts{}, &models.HistoryClientsDebt{}, &models.OverdueDebt{}, &models.BlackList{}, &models.Income{}, &models.Outcome{}, &models.Product{}, &models.ProductHistory{})
	DB.Model(&models.Debts{}).AddForeignKey("id_client", "clients(id)", "RESTRICT", "RESTRICT")
	DB.Model(&models.HistoryClientsDebt{}).AddForeignKey("id_client", "clients(id)", "RESTRICT", "RESTRICT")
	DB.Model(&models.OverdueDebt{}).AddForeignKey("id_client", "clients(id)", "RESTRICT", "RESTRICT")
	DB.Model(&models.BlackList{}).AddForeignKey("id_client", "clients(id)", "RESTRICT", "RESTRICT")
	DB.Model(&models.Income{}).AddForeignKey("id_client", "clients(id)", "RESTRICT", "RESTRICT")
	DB.Model(&models.Outcome{}).AddForeignKey("id_client", "clients(id)", "RESTRICT", "RESTRICT")
	DB.Model(&models.ProductHistory{}).AddForeignKey("id_client", "clients(id)", "RESTRICT", "RESTRICT")
	DB.Model(&models.ProductHistory{}).AddForeignKey("id_product", "products(id)", "RESTRICT", "RESTRICT")
	return DB
}
