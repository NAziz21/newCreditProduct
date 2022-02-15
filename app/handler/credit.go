package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/NAziz21/newCreditProduct/pkg/models"
	"github.com/NAziz21/newCreditProduct/settings/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Выдача кредита Post
func IssuancePost(ctx *gin.Context) {
	var request models.BodyRequest
	var dataFromDB models.Clients

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

	// проверка сущестоваоние клиента по номеру и получение его ID
	if err := database.DB.Raw("SELECT * from clients WHERE phone = ?", request.Phone).Scan(&dataFromDB).Error; err != nil {
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

	// Проверка долга у клиента
	totalDebt, err := TotalDebt(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Total Debt Error:", err)
		return
	}

	// Проверка просроченного долга
	overdueDebt, err := OverdueDebt(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Overdue Debt Error:", err)
		return
	}

	// Получение суммы дохода клиента
	financial, err := IncomeOutcome(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Cannot get income from DB")
		return
	}

	// Response
	responseStruct := models.IssuanceResponse{
		Message:     "Credit has been proved!",
		Status:      200,
		Name:        dataFromDB.Name,
		Surname:     dataFromDB.Surname,
		Debt:        totalDebt.Debt,
		OverDueDebt: overdueDebt,
		Income:      financial.Income,
		Outcome:     financial.Outcome,
	}

	totalIncome := responseStruct.Income - responseStruct.Outcome
	quantity := totalIncome / request.Amount

	if totalIncome <= 0 || responseStruct.Debt != 0 || responseStruct.OverDueDebt != 0 || quantity < 4 {
		responseStruct.Message = "Credit has not been proved!"
		ctx.JSON(http.StatusOK, gin.H{"CLIENT": responseStruct})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"CLIENT": responseStruct})
}

// Выдача кредита Get
func IssuanceGet(ctx *gin.Context) {
	var request models.HeaderRequest
	var dataFromDB models.Clients

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

	// проверка сущестоваоние клиента по номеру и получение его ID
	if err := database.DB.Raw("SELECT * from clients WHERE phone = ?", request.Phone).Scan(&dataFromDB).Error; err != nil {
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

	// Проверка долга у клиента
	totalDebt, err := TotalDebt(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Total Debt Error:", err)
		return
	}

	// Проверка просроченного долга
	overdueDebt, err := OverdueDebt(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Overdue Debt Error:", err)
		return
	}

	// Получение суммы дохода клиента
	financial, err := IncomeOutcome(ctx, dataFromDB.ID)
	if err != nil {
		log.Println("Cannot get income from DB")
		return
	}

	// Response
	responseStruct := models.IssuanceResponse{
		Message:     "Credit has been proved!",
		Status:      200,
		Name:        dataFromDB.Name,
		Surname:     dataFromDB.Surname,
		Debt:        totalDebt.Debt,
		OverDueDebt: overdueDebt,
		Income:      financial.Income,
		Outcome:     financial.Outcome,
	}

	totalIncome := responseStruct.Income - responseStruct.Outcome
	quantity := totalIncome / request.Amount

	if totalIncome <= 0 || responseStruct.Debt != 0 || responseStruct.OverDueDebt != 0 || quantity < 4 {
		responseStruct.Message = "Credit has not been proved!"
		ctx.JSON(http.StatusOK, gin.H{"CLIENT": responseStruct})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"CLIENT": responseStruct})
}

// Валидация Post Запроса клиента
func ValidationPostRequest(ctx *gin.Context, request models.BodyRequest) error {
	if request.Phone == "" || len(request.Phone) != 13 {
		log.Println("Phone number is too short")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Phone number is too short!",
			"Status":  400,
		})
		return errors.New("phone number is too short")
	}

	if request.Amount == 0 {
		log.Println("Amount is 0")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Amount is 0!",
			"Status":  400,
		})
		return errors.New("amount is 0")
	}
	return nil
}

// Валидация Get Запроса клиента
func ValidationGetRequest(ctx *gin.Context, request models.HeaderRequest) error {
	if request.Phone == "" || len(request.Phone) != 13 {
		log.Println("Phone number is too short")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Phone number is too short!",
			"Status":  400,
		})
		return errors.New("phone number is too short")
	}

	if request.Amount == 0 {
		log.Println("Amount is 0")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Amount is 0!",
			"Status":  400,
		})
		return errors.New("amount is 0")
	}
	return nil
}

// Сумма его доходов(входящие его суммы)
func IncomeOutcome(ctx *gin.Context, clientID int64) (*IncomeTotal, error) {
	var clientsFS IncomeTotal

	// получение дохода клиента
	if err := database.DB.Table("incomes").Select("SUM(salary + rent + extra_income) as income").Where("id_client = ?", clientID).First(&clientsFS).Error; err != nil {
		log.Println("Error:", err)
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{
			"Message": "Cannot get a client information from DB!",
			"Status":  204,
		})
		return nil, err
	}

	// получение расходов клиента
	if err := database.DB.Table("outcomes").Select("SUM(monthly_spendings + public_services + tax) as outcome").Where("id_client = ?", clientID).First(&clientsFS).Error; err != nil {
		log.Println("Error:", err)
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{
			"Message": "Cannot get a client information from DB!",
			"Status":  204,
		})
		return nil, err
	}
	return &clientsFS, nil
}
