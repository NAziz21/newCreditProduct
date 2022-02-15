package models

import "time"

//Таблица создания Клиента
type Clients struct {
	ID       int64     `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	Name     string    `gorm:"column:name" json:"name"`
	Surname  string    `gorm:"column:surname" json:"surname"`
	Phone    string    `gorm:"column:phone; unique" json:"phone"`
	Passport string    `gorm:"column:passport; unique" json:"passport"`
	Active   string    `gorm:"column:active;default:'active'" json:"active"`
	Created  time.Time `gorm:"autoCreateTime" json:"created"`
	Updated  time.Time `gorm:"autoUpdateTime" json:"updated"`
}

//Таблица Доходов и расходов клиента
type Income struct {
	ID         int64     `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	ClientID   int64     `gorm:"column:id_client" json:"id_client"`
	Salary     int64     `gorm:"column:salary;default:0" json:"salary"`
	Automobile string    `gorm:"column:automobile;default:'no'" json:"automobile"`
	Property   string    `gorm:"column:property;default:'no'" json:"property"`
	Extra      int64     `gorm:"column:extra_income;default:0" json:"extra_income"`
	Rent       int64     `gorm:"column:rent;default:0" json:"rent"`
	Created    time.Time `gorm:"autoCreateTime" json:"created"`
	Updated    time.Time `gorm:"autoUpdateTime" json:"updated"`
}

type Outcome struct {
	ID               int64     `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	ClientID         int64     `gorm:"column:id_client" json:"id_client"`
	MonthlySpendings int64     `gorm:"column:monthly_spendings;default:0" json:"monthly_spendings"`
	PublicService    int64     `gorm:"column:public_services;default:0" json:"public_services"`
	Tax              int64     `gorm:"column:tax;default:0" json:"tax"`
	Created          time.Time `gorm:"autoCreateTime" json:"created"`
	Updated          time.Time `gorm:"autoUpdateTime" json:"updated"`
}

type IssuanceResponse struct {
	Message     string `jspn:"Message"`
	Status      int64  `jspn:"Status"`
	Name        string `jspn:"Name"`
	Surname     string `jspn:"Surname"`
	Debt        int64  `jspn:"Debt"`
	OverDueDebt int64  `jspn:"Overdue debt"`
	Income      int64  `jspn:"Income"`
	Outcome     int64  `jspn:"Outcome"`
}

type TotalIncome struct {
	Income int64 `gorm:"income"`
}

//Таблица для Долга Клиента
type Debts struct {
	ID              int64     `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	ClientID        int64     `gorm:"column:id_client" json:"id_client"`
	AmountOfDebt    int64     `gorm:"column:debt_amount;default:100" json:"debt_amount"`
	Active          bool      `gorm:"column:active;default:true" json:"active"`
	StartDateOfDebt time.Time `gorm:"column:start_date" json:"start_date"`
	EndDateOfDebt   time.Time `gorm:"column:end_date" json:"end_date"`
}

type DebtResponse struct {
	Message     string `jspn:"Message"`
	Status      int64  `jspn:"Status"`
	Debt        int64  `jspn:"Debt"`
	PaidDebt    int64  `jspn:"Paid debt"`
	OverDueDebt int64  `jspn:"Overdue debt"`
}

// Таблица для историй клиента
type HistoryClientsDebt struct {
	ID              int64     `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	ClientID        int64     `gorm:"column:id_client" json:"id_client"`
	AmountOfDebt    int64     `gorm:"column:debt_amount;default:100" json:"debt_amount"`
	LoanRepayment   string    `gorm:"column:loan_repayment" json:"loan_repayment"`
	StartDateOfDebt time.Time `gorm:"column:start_date" json:"start_date"`
	EndDateOfDebt   time.Time `gorm:"column:end_date" json:"end_date"`
}

// Таблица для просроченного долга
type OverdueDebt struct {
	ID              int64     `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	ClientID        int64     `gorm:"column:id_client" json:"id_client"`
	AmountOfDebt    int64     `gorm:"column:debt_amount;default:100" json:"debt_amount"`
	StartDateOfDebt time.Time `gorm:"column:start_date" json:"start_date"`
	EndDateOfDebt   time.Time `gorm:"column:end_date" json:"end_date"`
}

// Запрос клиента
type BodyRequest struct {
	Phone   string `json:"phone"`
	Product string `form:"product"`
	Amount  int64  `json:"amount"`
}

// Запрос клиента header
type HeaderRequest struct {
	Phone   string `form:"phone"`
	Product string `form:"product"`
	Amount  int64  `form:"amount"`
}

// Получение ID клиента c Header
type DebtCheck struct {
	ClientID int64 `form:"clientID"`
}

// Для получения общей суммы долга и погашенного
type TotalDebt struct {
	Debt   int64 `gorm:"column:debt"`
	PayOff int64 `gorm:"column:pay_off"`
}

type BlackList struct {
	ID       int64     `gorm:"column:id;primary_key;autoIncrement"`
	IDClient int64     `gorm:"column:id_client"`
	Active   string    `gorm:"column:active;default:'active'" json:"active"`
	Date     time.Time `gorm:"date"`
}

// Таблица кредитного продукта
type Product struct {
	ID         int64     `gorm:"column:id;primary_key;autoIncrement"`
	Name       string    `gorm:"column:name" json:"name"`
	Amount     int64     `gorm:"column:amount" json:"amount"`
	Prepayment int64     `gorm:"column:prepayment" json:"prepayment"`
	AnnualRate int64     `gorm:"column:annual_rate" json:"annual rate"`
	StartDate  time.Time `gorm:"columns:start_date;default:'2021-01-01'" json:"start_date"`
	EndaDate   time.Time `gorm:"column:end_date;default:'2021-01-01'" json:"end_date"`
}

// Таблица кредитного продукта
type ProductPurchase struct {
	ID         int64  `gorm:"column:id;primary_key;autoIncrement"`
	Name       string `gorm:"column:name" json:"name"`
	Amount     int64  `gorm:"column:amount" json:"amount"`
	Prepayment int64  `gorm:"column:prepayment;default:0" json:"prepayment"`
	AnnualRate int64  `gorm:"column:annual_rate" json:"annual rate"`
	StartDate  time.Time
	StartPars  string `gorm:"columns:start_date" json:"start_date"`
	EndaDate   time.Time
	EndParse   string `gorm:"column:end_date" json:"end_date"`
}

// Считывание Header
type ProductResponseHeader struct {
	Name       string `form:"name"`
	Amount     int64  `form:"amount"`
	Prepayment int64  `form:"prepayment"` // Проценты
	AnnualRate int64  `form:"annualRate"` // Проценты
}

// Считывание Body
type ProductResponseBody struct {
	Name       string `json:"name"`
	Amount     int64  `json:"amount"`
	Prepayment int64  `json:"prepayment"` // Проценты
	AnnualRate int64  `json:"annualRate"` // проценты
}

// Таблица кредитного продукта
type ProductHistory struct {
	ID          int64     `gorm:"column:id;primary_key;autoIncrement"`
	ClientID    int64     `gorm:"column:id_client" json:"id_client"`
	ProductID   int64     `gorm:"column:id_product" json:"id_product"`
	ProductName string    `gorm:"column:product" json:"product"`
	StartDate   time.Time `gorm:"columns:start_date" json:"start_date"`	
	EndaDate    time.Time `gorm:"columns:end_date" json:"end_date"`
}
