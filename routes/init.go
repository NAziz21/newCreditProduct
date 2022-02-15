package routes

import (
	"github.com/NAziz21/newCreditProduct/app/handler"
	"github.com/NAziz21/newCreditProduct/settings"
	"github.com/NAziz21/newCreditProduct/settings/database"
	"github.com/gin-gonic/gin"
)

func Init() {
	settings.ConfigSetup("config/settings.json")

	database.ConnectDB()

	r := Router()

	r.Run(":" + settings.Config.PortRun)
}

func Router() *gin.Engine {
	router := gin.Default()
	subrouter := router.Group("/api/v1")

	// 1. Выдача кредита
	credit := subrouter.Group("/credit")
	{
		credit.POST("/issuance", handler.IssuancePost) // выдача кредита Post
		credit.GET("/issuance", handler.IssuanceGet)   // выдача кредита Get
	}

	// 2. Проверка существубщих задолженностей клиента
	subrouter.GET("debt", handler.DebtGet)       // Проверка его существующих задолженностей
	subrouter.POST("debt", handler.DebtPost)     // Проверка его существующих задолженностей
	subrouter.GET("debt/list", handler.DebtList) // Проверка его существующих задолженностей Клиента

	// 3. Создание кредитного продукта
	subrouter.GET("/product", handler.CreateProductGet)
	subrouter.POST("/product", handler.CreateProductPost)

	// 4. Создание скоринг системы
	subrouter.GET("/scoring", handler.ScoringGet)
	subrouter.POST("/scoring", handler.ScoringPost)

	// 5. Создание коллекторской части
	collector := subrouter.Group("/collector")
	{
		collector.GET("/debts", handler.CollectorDebt)           // Проверка его существующих задолженностей всех клиентов и подготовка коллекторской части
		collector.GET("/overdues", handler.CollectorOverdueDebt) // Проверка его существующих задолженностей всех клиентов и подготовка коллекторской части
	}

	// 6. Клиент подключает кредитный продукт
	subrouter.POST("/client/product", handler.Product)

	// Клиентская часть: Добавление, изменение и т.д. 
	subrouter.POST("/client", handler.AddClientH)                             // Добавление клиента
	subrouter.PATCH("/client", handler.UpdateClientH)                         // Изменение клиента
	subrouter.POST("/client/get", handler.GetClientH)                         // Получение клиента
	subrouter.GET("/client/all", handler.GetAllClientsH)                      // Получение всех клиентов
	subrouter.PATCH("/client/delete", handler.DeleteClient)                   // Удаление клиента(меняем его активность только)
	subrouter.POST("/client/financeStatement", handler.AddFinancialStatement) // Добавление Доходы/Расходы клиента
	subrouter.GET("/client/income", handler.IncomeA)                          // Получение его сумму входящих сумм

	return router
}
