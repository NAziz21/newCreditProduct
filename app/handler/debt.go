package handler

import (
	"log"
	"net/http"

	"github.com/NAziz21/newCreditProduct/app/service"
	"github.com/NAziz21/newCreditProduct/pkg/models"
	"github.com/NAziz21/newCreditProduct/settings/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Проверка долга клиента Get
func DebtGet(ctx *gin.Context) {
	var request models.HeaderRequest
	var client models.Clients

	if err := ctx.BindHeader(&request); err != nil {
		log.Println("Cannot read a client's request:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot read a header request",
			"Status":  400,
		})
		return
	}

	// Validation
	if err := ValidationGetRequest(ctx, request); err != nil {
		log.Println("Validation failed")
		return
	}

	// проверка сущестоваоние клиента по номеру и получение его ID
	if err := database.DB.Raw("SELECT * from clients WHERE phone = ?", request.Phone).Scan(&client).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("Error:", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Message": "This client is a new one!",
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

	// Общая сумма долга и погашения
	totalDebt, err := TotalDebt(ctx, client.ID)
	if err != nil {
		log.Println("Total Debt Error:", err)
		return
	}

	// просроченный долг
	overdueDebt, err := OverdueDebt(ctx, client.ID)
	if err != nil {
		log.Println("Overdue Debt Error:", err)
		return
	}

	response := models.DebtResponse{
		Message:     "Done!",
		Status:      200,
		Debt:        totalDebt.Debt,
		PaidDebt:    totalDebt.PayOff,
		OverDueDebt: overdueDebt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Client's Info": response,
	})
}

// Проверка долга клиента Get
func DebtPost(ctx *gin.Context) {
	var request models.BodyRequest
	var client models.Clients

	// получение запроса
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Print("Internal Error:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "Cannot read client's request!",
		})
		return
	}

	// Validation
	if err := ValidationPostRequest(ctx, request); err != nil {
		log.Println("Validation failed")
		return
	}

	// проверка сущестоваоние клиента по номеру и получение его ID
	if err := database.DB.Raw("SELECT * from clients WHERE phone = ?", request.Phone).Scan(&client).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("Error:", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Message": "This client is a new one!",
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

	// Общая сумма долга и погашения
	totalDebt, err := TotalDebt(ctx, client.ID)
	if err != nil {
		log.Println("Total Debt Error:", err)
		return
	}

	// Общая сумма долга и погашения
	overdueDebt, err := OverdueDebt(ctx, client.ID)
	if err != nil {
		log.Println("Overdue Debt Error:", err)
		return
	}

	response := models.DebtResponse{
		Message:     "Done!",
		Status:      200,
		Debt:        totalDebt.Debt,
		PaidDebt:    totalDebt.PayOff,
		OverDueDebt: overdueDebt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Client's Info": response,
	})
}

// Общая сумма долга и погашения
func TotalDebt(ctx *gin.Context, clientID int64) (*models.TotalDebt, error) {
	response, err := service.ClientsDebtTotalS(clientID)
	if err != nil {
		log.Println("Internal error(Clients Debt Check Total):", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal error",
			"Status":  500,
		})
		return nil, err
	}
	return response, nil
}

// Общая сумма долга и погашения
func OverdueDebt(ctx *gin.Context, clientID int64) (int64, error) {
	response, err := service.CheckOverdueDebtS(clientID)
	if err != nil {
		log.Println("Internal error(Clients Overdue Total):", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal error",
			"Status":  500,
		})
		return 0, err
	}
	return response, nil
}

// Список долгов клиента
func DebtList(ctx *gin.Context) {
	var client models.Clients
	var request models.HeaderRequest

	if err := ctx.BindHeader(&request); err != nil {
		log.Println("Cannot read a client's request:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot read a header request",
			"Status":  400,
		})
		return
	}

	// проверка сущестоваоние клиента по номеру и получение его ID
	if err := database.DB.Raw("SELECT * from clients WHERE phone = ?", request.Phone).Scan(&client).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("Error:", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Message": "This client is a new one!",
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
	response, err := service.ClientsDebtCheckS(client.ID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("such phone doesn't exists")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Message": "Such client doesn't exists!",
				"Status":  500,
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal error",
			"Status":  500,
		})
		log.Println("Internal error(Clients Debt Check):", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Message":       200,
		"Client's Debt": response,
	})
}
