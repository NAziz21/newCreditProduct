package service

import (
	"log"

	"github.com/NAziz21/newCreditProduct/pkg/models"
	"github.com/NAziz21/newCreditProduct/settings/database"
	"github.com/jinzhu/gorm"
)

// Получение списка кредитов
func ClientsDebtCheckS(clientID int64) ([]*models.DebtResponse, error) {
	listOfDebts := make([]*models.DebtResponse, 0)

	if err := database.DB.Table("debts").Where("id_client = ? AND active = 'true'", clientID).Find(&listOfDebts).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("such client has not a debt!")
			return nil, err
		}
		log.Println("Cannot take debts from db. Error:", err)
		return nil, err
	}

	return listOfDebts, nil
}

// Общий долг и погашение
func ClientsDebtTotalS(clientID int64) (*models.TotalDebt, error) {
	var ClientInfo models.TotalDebt

	// Получение общего сумма долга
	if err := database.DB.Raw("SELECT SUM(debt_amount) as debt FROM debts	WHERE id_client = ? AND active = 'true'", clientID).First(&ClientInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("such client doesn't exist in DB!")
			return nil, err
		}
		log.Println("Cannot do sum a client's debt")
		return nil, err
	}

	// Получение общей суммы погашенного долга
	if err := database.DB.Raw("SELECT SUM(debt_amount) as pay_off FROM debts	WHERE id_client = ? AND active = 'false'", clientID).First(&ClientInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("such client doesn't exist in DB!")
			return nil, err
		}
		log.Println("Cannot do sum a client's payed debts")
		return nil, err
	}

	return &ClientInfo, nil
}

type OverdueDebt struct {
	Overdue int64 `gorm:"column:overdue"`
}

// Проверка просроченного долга у клиента
func CheckOverdueDebtS(clientID int64) (int64, error) {
	var clientInfoFromDB OverdueDebt

	if err := database.DB.Raw("SELECT SUM(debt_amount) as overdue FROM overdue_debts	WHERE id_client = ?", clientID).First(&clientInfoFromDB).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("such client doesn't exist in DB!")
			return 0, err
		}
		log.Println("Cannot do sum a client's payed debts")
		return 0, err
	}

	return clientInfoFromDB.Overdue, nil
}
