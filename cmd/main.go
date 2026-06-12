package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"accounting/internal/handler"
	"accounting/internal/repository"
	"accounting/internal/service"
)

func main() {
	godotenv.Load()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("veritabanına bağlanılamadı: %v", err)
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"json": func(v any) (template.JS, error) {
			b, err := json.Marshal(v)
			return template.JS(b), err
		},
	}).ParseGlob("web/templates/*.html")
	if err != nil {
		log.Fatalf("şablon hatası: %v", err)
	}

	accountRepo := repository.NewAccountRepository(db)
	accountSvc := service.NewAccountService(accountRepo)

	incomeCategoryRepo := repository.NewIncomeCategoryRepository(db)
	incomeCategorySvc := service.NewIncomeCategoryService(incomeCategoryRepo)

	incomeRepo := repository.NewIncomeRepository(db)
	incomeSvc := service.NewIncomeService(incomeRepo)

	expenseCategoryRepo := repository.NewExpenseCategoryRepository(db)
	expenseCategorySvc := service.NewExpenseCategoryService(expenseCategoryRepo)

	expenseRepo := repository.NewExpenseRepository(db)
	expenseSvc := service.NewExpenseService(expenseRepo)

	tourCategoryRepo := repository.NewTourCategoryRepository(db)
	tourCategorySvc := service.NewTourCategoryService(tourCategoryRepo)

	roomRepo := repository.NewRoomRepository(db)
	roomSvc := service.NewRoomService(roomRepo)

	tourRepo := repository.NewTourRepository(db)
	tourSvc := service.NewTourService(tourRepo)

	clientRepo := repository.NewClientRepository(db)
	clientSvc := service.NewClientService(clientRepo)

	settingRepo := repository.NewSettingRepository(db)
	settingSvc := service.NewSettingService(settingRepo)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	discountCategoryRepo := repository.NewDiscountCategoryRepository(db)
	discountCategorySvc := service.NewDiscountCategoryService(discountCategoryRepo)

	discountRepo := repository.NewDiscountRepository(db)
	discountSvc := service.NewDiscountService(discountRepo)

	orderRepo := repository.NewOrderRepository(db)
	orderSvc := service.NewOrderService(orderRepo, incomeRepo)

	accountHandler := handler.NewAccountHandler(accountSvc)
	incomeCategoryHandler := handler.NewIncomeCategoryHandler(incomeCategorySvc)
	incomeHandler := handler.NewIncomeHandler(incomeSvc)
	expenseCategoryHandler := handler.NewExpenseCategoryHandler(expenseCategorySvc)
	expenseHandler := handler.NewExpenseHandler(expenseSvc)
	tourCategoryHandler := handler.NewTourCategoryHandler(tourCategorySvc)
	roomHandler := handler.NewRoomHandler(roomSvc)
	tourHandler := handler.NewTourHandler(tourSvc)
	clientHandler := handler.NewClientHandler(clientSvc)
	settingHandler := handler.NewSettingHandler(settingSvc)
	userHandler := handler.NewUserHandler(userSvc)
	discountCategoryHandler := handler.NewDiscountCategoryHandler(discountCategorySvc)
	discountHandler := handler.NewDiscountHandler(discountSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	pageHandler := handler.NewPageHandler(accountSvc, incomeCategorySvc, incomeSvc, expenseCategorySvc, expenseSvc, tourCategorySvc, roomSvc, tourSvc, clientSvc, settingSvc, userSvc, discountCategorySvc, discountSvc, orderSvc)

	router := handler.NewRouter(accountHandler, incomeCategoryHandler, incomeHandler, expenseCategoryHandler, expenseHandler, tourCategoryHandler, roomHandler, tourHandler, clientHandler, settingHandler, userHandler, discountCategoryHandler, discountHandler, orderHandler, pageHandler, tmpl)
	router.Run(":" + os.Getenv("PORT"))
}
