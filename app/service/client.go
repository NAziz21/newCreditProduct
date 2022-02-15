package service

import (
	"log"
	"strings"
	"time"

	"github.com/NAziz21/newCreditProduct/pkg/models"
	"github.com/NAziz21/newCreditProduct/settings/database"
)

// Добавление клиента
func AddClient(clientInfo models.Clients) (*models.Clients, error) {

	client := models.Clients{
		Name:     TittleString(clientInfo.Name),
		Surname:  TittleString(clientInfo.Surname),
		Phone:    "+" + clientInfo.Phone,
		Passport: clientInfo.Passport,
		Created:  time.Now().Local(),
		Updated:  time.Now().Local(),
	}

	if err := database.DB.Table("clients").Create(&client).Error; err != nil {
		log.Println("Cannot add a client's info")
		return nil, err
	}

	return &client, nil
}

// Добавление Financial Statement клиента
func AddFS(clientInfo models.Income) (*models.Income, error) {

	client := models.Income{
		ClientID:   clientInfo.ClientID,
		Salary:     clientInfo.Salary,
		Automobile: clientInfo.Automobile,
		Property:   clientInfo.Property,
		Extra:      clientInfo.Extra,
		Rent:       clientInfo.Rent,
		Created:    time.Now().Local(),
		Updated:    time.Now().Local(),
	}
	if err := database.DB.Table("financial_statements").Create(&client).Error; err != nil {
		log.Println("Cannot add a client's info")
		return nil, err
	}
	return &client, nil
}

// верхний регистр первой буквы
func TittleString(name string) (newName string) {
	lowerCase := strings.ToLower(name)
	newName = strings.Title(lowerCase)
	return newName
}

// Проверка клиент старый или новый
func ExistedClientS(clientInfoRequest models.BodyRequest) (*models.Clients, error) {
	var clientInfoFromDB models.Clients
	query := database.DB.Table("clients").Where("phone = ?", clientInfoRequest.Phone).First(&clientInfoFromDB)

	if err := query.Error; err != nil {
		log.Println("such client doesn't exist in DB!")
		return nil, err
	}
	return &clientInfoFromDB, nil
}
