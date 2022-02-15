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

/*
 Пример запроса
{
    "name": "Amal 2019",
    "amount": 50000,
    "prepayment": 0,
    "annualRate": 25
}
*/

// Создание продукта Get
func CreateProductGet(ctx *gin.Context) {
	var request models.ProductResponseHeader
	var oldProduct models.Product

	if err := ctx.BindHeader(&request); err != nil {
		log.Println("Cannot read a client's request:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot read a header request",
			"Status":  400,
		})
		return
	}

	// Validation
	if err := ProductValidationGetRequest(ctx, request); err != nil {
		log.Println("Validation failed")
		return
	}

	// Проверка на уникальность
	if err := database.DB.Raw("SELECT * from products WHERE name = ?", request.Name).Scan(&oldProduct).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// Создание нового кредитного продукта
			newProduct := models.ProductResponseHeader{
				Name:       request.Name,
				Amount:     request.Amount,
				AnnualRate: request.AnnualRate,
				Prepayment: request.Prepayment,
			}

			if err := database.DB.Table("products").Create(&newProduct).Error; err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"Message": "Cannot add a client's info!",
					"Error":   err,
					"Status":  400,
				})
				log.Println("Cannot add a client's info")
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"Message":      "Done!",
				"Status":       200,
				"Product Name": request.Name,
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

	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"Message": "Such Product exists in DB!",
		"Status":  400,
	})
}

// Создание продукта Post
func CreateProductPost(ctx *gin.Context) {
	var request models.ProductResponseBody
	var oldProduct models.Product

	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Println("Cannot read a client's request:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Cannot read a header request",
			"Status":  400,
		})
		return
	}

	// Validation
	if err := ProductValidationPost(ctx, request); err != nil {
		log.Println("Validation failed")
		return
	}

	// Проверка на уникальность
	if err := database.DB.Raw("SELECT * from products WHERE name = ?", request.Name).Scan(&oldProduct).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// Создание нового кредитного продукта
			newProduct := models.ProductResponseBody{
				Name:       request.Name,
				Amount:     request.Amount,
				AnnualRate: request.AnnualRate,
				Prepayment: request.Prepayment,
			}

			if err := database.DB.Table("products").Create(&newProduct).Error; err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"Message": "Cannot add a client's info!",
					"Error":   err,
					"Status":  400,
				})
				log.Println("Cannot add a client's info")
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"Message":      "Done!",
				"Status":       200,
				"Product Name": request.Name,
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

	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"Message": "Such Product exists in DB!",
		"Status":  400,
	})
}

// Валидация Get Запроса клиента
func ProductValidationGetRequest(ctx *gin.Context, request models.ProductResponseHeader) error {
	if request.Name == "" || len(request.Name) <= 3 {
		log.Println("Product name are too short")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Product name  are too short!",
			"Status":  400,
		})
		return errors.New("product name  are too short")
	}

	if request.Amount == 0 || request.AnnualRate == 0 {
		log.Println("Amount or annual are 0")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Amount or Prepayment or annual are 0!",
			"Status":  400,
		})
		return errors.New("amount or annual are 0")
	}
	return nil
}

// Валидация Get Запроса клиента
func ProductValidationPost(ctx *gin.Context, request models.ProductResponseBody) error {
	if request.Name == "" || len(request.Name) <= 3 {
		log.Println("Product name are too short")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Product name  are too short!",
			"Status":  400,
		})
		return errors.New("product name  are too short")
	}

	if request.Amount == 0 || request.AnnualRate == 0 {
		log.Println("Amount or annual are 0")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Message": "Amount or Prepayment or annual are 0!",
			"Status":  400,
		})
		return errors.New("amount or annual are 0")
	}
	return nil
}
