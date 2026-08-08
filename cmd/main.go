package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"accounting/internal/auth"
	"accounting/internal/googlesheets"
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
		"now": time.Now,
		"ddmmyyyy": func(s string) string {
			if len(s) < 10 {
				return s
			}
			return s[8:10] + "." + s[5:7] + "." + s[0:4]
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

	flightRepo := repository.NewFlightRepository(db)
	flightSvc := service.NewFlightService(flightRepo)

	tourRepo := repository.NewTourRepository(db)
	tourSvc := service.NewTourService(tourRepo)

	clientRepo := repository.NewClientRepository(db)
	clientSvc := service.NewClientService(clientRepo)

	settingRepo := repository.NewSettingRepository(db)
	settingSvc := service.NewSettingService(settingRepo)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	sessions := auth.NewStore()

	discountCategoryRepo := repository.NewDiscountCategoryRepository(db)
	discountCategorySvc := service.NewDiscountCategoryService(discountCategoryRepo)

	discountRepo := repository.NewDiscountRepository(db)
	discountSvc := service.NewDiscountService(discountRepo)

	orderRepo := repository.NewOrderRepository(db)
	orderSvc := service.NewOrderService(orderRepo, incomeRepo, discountRepo, tourRepo)

	taskRepo := repository.NewTaskRepository(db)
	taskSvc := service.NewTaskService(taskRepo)

	sheetsClient, err := googlesheets.NewClient(context.Background(), googlesheets.CredentialsPath())
	if err != nil {
		log.Printf("Google Sheets qoşulmadı: %v", err)
	}

	sheetLinkRepo := repository.NewSheetLinkRepository(db)
	sheetLinkSvc := service.NewSheetLinkService(sheetLinkRepo)

	accountHandler := handler.NewAccountHandler(accountSvc, incomeSvc, expenseSvc, orderSvc)
	incomeCategoryHandler := handler.NewIncomeCategoryHandler(incomeCategorySvc)
	incomeHandler := handler.NewIncomeHandler(incomeSvc)
	expenseCategoryHandler := handler.NewExpenseCategoryHandler(expenseCategorySvc)
	expenseHandler := handler.NewExpenseHandler(expenseSvc)
	tourCategoryHandler := handler.NewTourCategoryHandler(tourCategorySvc)
	roomHandler := handler.NewRoomHandler(roomSvc)
	flightHandler := handler.NewFlightHandler(flightSvc)
	tourHandler := handler.NewTourHandler(tourSvc)
	clientHandler := handler.NewClientHandler(clientSvc)
	settingHandler := handler.NewSettingHandler(settingSvc)
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(userSvc, sessions)
	discountCategoryHandler := handler.NewDiscountCategoryHandler(discountCategorySvc)
	discountHandler := handler.NewDiscountHandler(discountSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	taskHandler := handler.NewTaskHandler(taskSvc)
	pageHandler := handler.NewPageHandler(accountSvc, incomeCategorySvc, incomeSvc, expenseCategorySvc, expenseSvc, tourCategorySvc, roomSvc, flightSvc, tourSvc, clientSvc, settingSvc, userSvc, discountCategorySvc, discountSvc, orderSvc, taskSvc)
	sheetsImportHandler := handler.NewSheetsImportHandler(sheetsClient, sheetLinkSvc, tourSvc, clientSvc, orderSvc)

	router := handler.NewRouter(authHandler, sessions, accountHandler, incomeCategoryHandler, incomeHandler, expenseCategoryHandler, expenseHandler, tourCategoryHandler, roomHandler, flightHandler, tourHandler, clientHandler, settingHandler, userHandler, discountCategoryHandler, discountHandler, orderHandler, taskHandler, pageHandler, sheetsImportHandler, tmpl)
	router.Run(":" + os.Getenv("PORT"))
}
