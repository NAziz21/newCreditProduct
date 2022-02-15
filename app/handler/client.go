package handler

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/NAziz21/newCreditProduct/app/service"
	"github.com/NAziz21/newCreditProduct/pkg/models"
	"github.com/NAziz21/newCreditProduct/settings/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Добавление клиента
func AddClientH(ctx *gin.Context) {
	var clientInfo models.Clients

	if err := ctx.ShouldBindJSON(&clientInfo); err != nil {
		log.Print("Internal Error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Fields are empty!",
		})
		return
	}

	validation, err := ValidationClientRequest(clientInfo)
	if err != nil {
		log.Print("Internal error:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Error!",
			"Status":  500,
		})
		return
	}

	if !validation {
		log.Print("Validation Error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": err,
			"Status":  400,
		})
		return
	}

	// Проверка пользователя
	err = UniqueCheck("clients", "phone = ? OR passport = ?", clientInfo.Phone, clientInfo.Passport)
	if err != nil && err == gorm.ErrRecordNotFound {
		// Добавление пользователя
		response, err := service.AddClient(clientInfo)
		if err != nil {
			log.Println("Error:", err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Message":     "Client was successfuly added!",
			"Status":      200,
			"Client Name": response.Name,
		})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{
		"Message": "Such client is already existed!",
		"Status":  400,
	})
}

// Проверка Существующего клиента
func GetClientH(ctx *gin.Context) {
	var clientRequest models.BodyRequest

	if err := ctx.ShouldBindJSON(&clientRequest); err != nil {
		log.Println("Cannot read a client's request:", err)
		return
	}

	response, err := service.ExistedClientS(clientRequest)
	if err != nil {
		ctx.JSON(200, gin.H{
			"Message": "Such client doesn't exists in DB!",
			"Status":  200,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"Message":        "Such client exists in DB!",
		"Status":         200,
		"Client Name":    response.Name,
		"Client Surname": response.Surname,
	})
}

// Получение всех фильмов
func GetAllClientsH(ctx *gin.Context) {
	var getAll []models.Clients

	if err := database.DB.Table("clients").Find(&getAll).Error; err != nil {
		log.Println("Internal Error(Get All Client):", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Status":  400,
			"Message": err,
		})
	}

	ctx.JSON(200, gin.H{
		"Status":  400,
		"Clients": getAll,
	})
}

// Обновление клиента
func UpdateClientH(ctx *gin.Context) {
	var updateRequest, dbInfo models.Clients

	// Обрабатываем Body Change
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		log.Println("Cannot read a request body:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot read a request body",
			"Status":  400,
		})
		return
	}

	// Проверка на существование ID
	if err := database.DB.Table("clients").Where("id = ?", updateRequest.ID).First(&dbInfo).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Message": "Such client doesn't exist",
				"Status":  400,
			})
			return
		}

		log.Println("Internal Error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot read a request body",
			"Status":  400,
		})
		return
	}

	// проверяем на уникальность
	uniqueCheckSupport, err := UniqueCheck2(updateRequest)
	if err != nil {
		log.Println("Internal Error(Update Client):", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Status": 400,
		})
		return
	}
	if !uniqueCheckSupport {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Such passport or phone is already existed!",
			"Status":  400,
		})
		return
	}

	// обновляем
	if err := database.DB.Model(&dbInfo).Update(updateRequest).Error; err != nil {
		log.Println("Cannot update:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot update",
			"Status":  400,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Done!"})
}

// удаление клиента путем измененение его активности (через body)
func DeleteClient(ctx *gin.Context) {
	var deleteRequest, infoFromDB models.Clients

	if err := ctx.ShouldBindJSON(&deleteRequest); err != nil {
		log.Println("Cannot read a request body:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot read a request body",
			"Status":  400,
		})
		return
	}

	// Проверка на существование ID
	if err := database.DB.Table("clients").Where("id = ?", deleteRequest.ID).First(&infoFromDB).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Message": "Such client doesn't exist",
				"Status":  400,
			})
			return
		}

		log.Println("Internal Error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "error",
			"Status":  400,
		})
		return
	}

	// обновляем
	if err := database.DB.Model(&infoFromDB).Update(deleteRequest).Error; err != nil {
		log.Println("Cannot delete:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot update",
			"Status":  400,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Message": "Done!",
		"Status":  400,
	})
}

// Валидация запроса клиента
func ValidationClientRequest(clientRequest models.Clients) (bool, error) {

	if len(clientRequest.Name) < 4 {
		log.Println("Client's Name is too short!")
		return false, errors.New("client's Name is too short")
	}

	if len(clientRequest.Surname) < 2 {
		log.Println("Client's surname is too short!")
		return false, errors.New("client's surname is too short")
	}

	re := regexp.MustCompile(`([0-9])`)
	submatch := re.FindAllString(clientRequest.Phone, -1)
	if len(submatch) != 12 {
		log.Println(submatch[1])
		log.Println("Incorrect number!")
		return false, errors.New("phone number is incorrect")
	}
	return true, nil
}

// Проверяем паспорт и телефон на уникальность
func UniqueCheck(nameOfTableDB string, columnNameOfTableQuery string, phone string, passport string) error {
	var clientInfoFromDB models.Clients

	query := database.DB.Table(nameOfTableDB)
	if err := query.Where(columnNameOfTableQuery, phone, passport).First(&clientInfoFromDB).Error; err != nil {
		log.Println("such login doesn't exists:", err)
		return err
	}
	return nil
}

// Проверяем логин и телефон на уникальность
func UniqueCheck2(usersInfo models.Clients) (bool, error) {
	var infoFromDB models.Clients

	query := database.DB.Table("clients")
	if err := query.Where("passport = ? OR phone = ?", infoFromDB.Passport, infoFromDB.Phone).First(&infoFromDB).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("such phone or passport doesn't exists in DB")
			return false, nil
		}
		log.Print("Internal error(Unique Check):", err)
		return false, err
	}
	return true, nil
}

// Добавление клиента
func AddFinancialStatement(ctx *gin.Context) {
	var clientInfo models.Income

	if err := ctx.ShouldBindJSON(&clientInfo); err != nil {
		log.Println("Cannot read a request body:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot read a request body",
			"Status":  400,
		})
		return
	}

	// Валидация запроса пользователя
	err := ValidationFinanceStatement(clientInfo)
	if err != nil {
		log.Print("Validation failed:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Validation failed!",
			"Error":   err,
			"Status":  400,
		})
		return
	}

	// Проверка не повторение id
	uniqueCheckID, err := UniqueCheckFinanceStatement(clientInfo)
	if err != nil {
		log.Println("Internal Error(Update Client):", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Status": 400,
		})
		return
	}
	if !uniqueCheckID {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Such id is already existed!",
			"Status":  400,
		})
		return
	}

	log.Println(uniqueCheckID)
	// Добавление в БД
	_, err = service.AddFS(clientInfo)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot add a client's finance",
			"Status":  400,
		})
		return
	}

	ctx.AbortWithStatusJSON(200, gin.H{
		"Message": "Done!",
		"Status":  200,
	})
}

// Проверяем логин и телефон на уникальность
func UniqueCheckFinanceStatement(usersInfo models.Income) (bool, error) {
	var infoFromDB models.Income
	query := database.DB.Table("financial_statements")
	if err := query.Where("id_client = ?", usersInfo.ClientID).First(&infoFromDB).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("such id doesn't exists in DB")
			return true, nil
		}
		log.Print("Internal error(Unique Check):", err)
		return false, err
	}
	return false, nil
}

// Валидация запроса клиента
func ValidationFinanceStatement(clientIfnfo models.Income) error {
	if clientIfnfo.ClientID == 0 {
		log.Println("Client's ID is absent!")
		return errors.New("client's ID is absent")
	}

	check := []string{"yes", "no"}
	for _, val := range check {
		if val == clientIfnfo.Automobile && val == clientIfnfo.Property {
			return nil
		}
	}
	return nil
}

type PhoneNumber struct {
	Phone string `form:"phone"`
}

type IncomeTotal struct {
	Income  int64 `gorm:"income"`
	Outcome int64 `gorm:"outcome"`
}

func IncomeA(ctx *gin.Context) {
	var phone PhoneNumber
	var infoFromDB models.Clients
	var clientsFS IncomeTotal

	// получение номера клиента через заголовок
	if err := ctx.BindHeader(&phone); err != nil {
		log.Println("Cannot read a header request:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot read a header request",
			"Status":  400,
		})
		return
	}

	log.Println(phone.Phone)
	// проверка сущестоваоние клиента по номеру и получение его ID
	if err := database.DB.Raw("SELECT id from clients WHERE phone = ?", phone.Phone).Scan(&infoFromDB).Error; err != nil {
		log.Println("Error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Internal Error!",
			"Error":   err,
			"Status":  400,
		})
		return
	}

	// получение информации клиента
	if err := database.DB.Table("financial_statements").Select("SUM(salary + rent + extra_income) as income").Where("id_client = ?", infoFromDB.ID).First(&clientsFS).Error; err != nil {
		log.Println("Error:", err)
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{
			"Message": "Cannot get a client information from DB!",
			"Status":  204,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": 200,
		"Income": clientsFS.Income,
	})
}

// Клиент берет кредит(подключает кредитный продукт)
func Product(ctx *gin.Context) {
	var request models.BodyRequest

	// получение запроса
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Print("Internal Error:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "Cannot read client's request!",
		})
		return
	}

	scoring, err := ScoringSystemPost(ctx, request)
	if err != nil {
		log.Println("Error(Purchase product):", err)
		return
	}

	if scoring.Status == 1 {
		log.Println("You cannot get a credit")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "You cannot get a credit!",
			"Status": 400,
		})
		return
	}

	if scoring.Status == 2 {
		log.Println("Your case will be on consideration! Our agent will call you!")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Message": "Your case will be on consideration! Our agent will call you!",
			"Status": 400,
		})
		return
	}

	if scoring.Status >= 3 {
		err := PurchaseProduct(ctx, request)
		if  err != nil {
			log.Println("error", err)
			return
		}
	}
}

func PurchaseProduct(ctx *gin.Context, request models.BodyRequest) error {
	var product models.Product
	var client models.Clients
	if err := database.DB.Table("products").Where("name = ?", request.Product).Scan(&product).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("Such product name doesn't exists")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Message": "Such product name doesn't exists!",
				"Status": 400,
			})
			return err
		}
		log.Println("Internal Error:", err)
		return err
	}

	// проверка сущестоваоние клиента по номеру и получение его ID
	if err := database.DB.Raw("SELECT * from clients WHERE phone = ?", request.Phone).Scan(&client).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("Error:", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Message": "This client is a new one!",
				"Status":  400,
			})
			return err 
		}

		log.Println("Error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Internal Error!",
			"Error":   err,
			"Status":  400,
		})
		return err
	} 

	addProduct := models.ProductHistory{
		ClientID: client.ID,
		ProductID: product.ID,
		ProductName: product.Name,
		StartDate: time.Now().Local(),
		EndaDate: time.Now().Local().Add( 360 * 24 * time.Hour),
	}

	if err := database.DB.Table("product_histories").Create(&addProduct).Error; err != nil {
		log.Println("Error while inserting to table product histories")
		return err
	}

	ctx.JSON(200, gin.H{
		"Message": "Well Done!",
		"Status": 200,
	}) 
	return nil
}
