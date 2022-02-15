package handler

import (
	"log"
	"net/http"

	"github.com/NAziz21/newCreditProduct/settings/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CollectorDebtResponse struct {
	Name       string `json:"Name"`
	Surname    string `json:"Surname"`
	DebtAmount int64  `json:"Debt_amount"`
}

type CollectorOverdueResponse struct {
	Name       string `json:"Name"`
	Surname    string `json:"Surname"`
	DebtAmount int64  `json:"Debt_amount"`
}

func CollectorDebt(ctx *gin.Context) {
	var listOfDebts []CollectorDebtResponse
	// Берем список долгов

	query := "SELECT debts.debt_amount, clients.name, clients.surname FROM debts JOIN clients ON debts.id_client = clients.id WHERE debts.active = 'true'"
	if err := database.DB.Raw(query).Scan(&listOfDebts).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("Table is empty")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Message": "Table is empty!",
				"Status":  400,
			})
			return
		}

		log.Println("Error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Internal Error!",
			"Error":   err,
			"Status":  400,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Message":      "Done!",
		"Status":       200,
		"List of Debts": listOfDebts,
	})
}

func CollectorOverdueDebt(ctx *gin.Context) {
	
	var listOfOverdue []CollectorOverdueResponse
	
	query := "SELECT overdue_debts.debt_amount, clients.name, clients.surname FROM overdue_debts JOIN clients ON overdue_debts.id_client = clients.id"
	if err := database.DB.Raw(query).Scan(&listOfOverdue).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("Table is empty")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Message": "Table is empty!",
				"Status":  400,
			})
			return
		}

		log.Println("Error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Internal Error!",
			"Error":   err,
			"Status":  400,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Message":      "Done!",
		"Status":       200,
		"List of Debts": listOfOverdue,
	})
}
