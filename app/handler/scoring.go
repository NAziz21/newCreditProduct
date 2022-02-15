package handler

import (
	"log"
	"net/http"

	"github.com/NAziz21/newCreditProduct/pkg/models"
	"github.com/NAziz21/newCreditProduct/settings/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const CreditDeny = 1
const Consideration = 2
const Approved = 3

// Скоринг Get
func ScoringGet(ctx *gin.Context) {
	var request models.HeaderRequest

	// получение запроса
	if err := ctx.BindHeader(&request); err != nil {
		log.Print("Internal Error:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "Cannot read client's request!",
		})
		return
	}

	err := ValidationGetRequest(ctx, request)
	if err != nil {
		log.Println("Validation failed")
		return
	}

	response, err := ScoringSystemGet(ctx, request)
	if err != nil {
		log.Println("Error in Scoring system", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Client's Info": response,
	})
}

// Скоринг Post
func ScoringPost(ctx *gin.Context) {
	var request models.BodyRequest

	// получение запроса
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Print("Internal Error:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "Cannot read client's request!",
		})
		return
	}

	err := ValidationPostRequest(ctx, request)
	if err != nil {
		log.Println("Validation failed")
		return
	}

	response, err := ScoringSystemPost(ctx, request)
	if err != nil {
		log.Println("Error in Scoring system", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Client's Info": response,
	})
}

// Система счетов Get запроса
func ScoringSystemGet(ctx *gin.Context, request models.HeaderRequest) (*models.IssuanceResponse, error) {
	var client models.Income
	var dataFromDB models.Clients
	var message string
	var status int64
	score := 0

	// проверка сущестоваоние клиента по номеру и получение его ID
	if err := database.DB.Raw("SELECT * from clients WHERE phone = ?", request.Phone).Scan(&dataFromDB).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("Error:", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Message": "This client is a new one!",
				"Status":  400,
			})
			return nil, err
		}

		log.Println("Error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Internal Error!",
			"Error":   err,
			"Status":  400,
		})
		return nil, err
	}

	// Проверка долга у клиента
	totalDebt, err := TotalDebt(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Total Debt Error:", err)
		return nil, err
	}

	// Проверка просроченного долга
	overdueDebt, err := OverdueDebt(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Overdue Debt Error:", err)
		return nil, err
	}

	// Получение суммы дохода клиента
	financial, err := IncomeOutcome(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Cannot get income from DB")
		return nil, err
	}

	if totalDebt.Debt == 0 {
		score += 1
	}

	if totalDebt.PayOff == 0 {
		score += 1
	}

	if overdueDebt == 0 {
		score += 1
	} else {
		score -= 1
	}

	if err := database.DB.Table("incomes").Where("id_client = ?", dataFromDB.ID).Scan(&client).Error; err != nil {
		log.Println("Scoring system:", err)
		return nil, err
	}

	if client.Automobile == "yes" {
		score += 1
	} else {
		score -= 1
	}

	if client.Property == "yes" {
		score += 2
	} else {
		score -= 2
	}

	if score <= 2 {
		message = "Credit denied!"
		status = CreditDeny

	}

	if score == 3 {
		message = "Need more time"
		status = Consideration
	}

	if score >= 4 {
		message = "Credit has been approved"
		status = Approved
	}

	// Response
	responseStruct := models.IssuanceResponse{
		Message:     message,
		Status:      status,
		Name:        dataFromDB.Name,
		Surname:     dataFromDB.Surname,
		Debt:        totalDebt.Debt,
		OverDueDebt: overdueDebt,
		Income:      financial.Income,
		Outcome:     financial.Outcome,
	}
	return &responseStruct, nil
}

// Система счетов Post запроса
func ScoringSystemPost(ctx *gin.Context, request models.BodyRequest) (*models.IssuanceResponse, error) {
	var client models.Income
	var dataFromDB models.Clients
	var message string
	var status int64
	score := 0

	// проверка сущестоваоние клиента по номеру и получение его ID
	if err := database.DB.Raw("SELECT * from clients WHERE phone = ?", request.Phone).Scan(&dataFromDB).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("Error:", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Message": "This client is a new one!",
				"Status":  400,
			})
			return nil, err
		}

		log.Println("Error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Internal Error!",
			"Error":   err,
			"Status":  400,
		})
		return nil, err
	}

	// Проверка долга у клиента
	totalDebt, err := TotalDebt(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Total Debt Error:", err)
		return nil, err
	}

	// Проверка просроченного долга
	overdueDebt, err := OverdueDebt(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Overdue Debt Error:", err)
		return nil, err
	}

	// Получение суммы дохода клиента
	financial, err := IncomeOutcome(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Cannot get income from DB")
		return nil, err
	}

	if totalDebt.Debt == 0 {
		score += 1
	}

	if totalDebt.PayOff == 0 {
		score += 1
	}

	if overdueDebt == 0 {
		score += 1
	} else {
		score -= 1
	}

	if err := database.DB.Table("incomes").Where("id_client = ?", dataFromDB.ID).Scan(&client).Error; err != nil {
		log.Println("Scoring system:", err)
		return nil, err
	}

	if client.Automobile == "yes" {
		score += 1
	} else {
		score -= 1
	}

	if client.Property == "yes" {
		score += 2
	} else {
		score -= 2
	}

	if score <= 2 {
		message = "Credit denied!"
		status = CreditDeny

	}

	if score == 3 {
		message = "Need more time"
		status = Consideration
	}

	if score >= 4 {
		message = "Credit has been approved"
		status = Approved
	}

	// Response
	responseStruct := models.IssuanceResponse{
		Message:     message,
		Status:      status,
		Name:        dataFromDB.Name,
		Surname:     dataFromDB.Surname,
		Debt:        totalDebt.Debt,
		OverDueDebt: overdueDebt,
		Income:      financial.Income,
		Outcome:     financial.Outcome,
	}
	return &responseStruct, nil
}
