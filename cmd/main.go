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

// applyMigration035 догоняет migrations/035_drop_tours_room_id.sql, если она ещё
// не применена: старая колонка tours.room_id (NOT NULL, FK на rooms) заставляет
// MySQL подставлять room_id=0 при создании тура через GORM, что валит вставку
// нарушением fk_tour_room, — комнаты тура теперь живут в tour_rooms.
func applyMigration035(db *gorm.DB) error {
	var count int64
	err := db.Raw(`
		SELECT COUNT(*) FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'tours' AND COLUMN_NAME = 'room_id'
	`).Scan(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	log.Println("migrasyon 035 tətbiq olunur: tours.room_id sütunu silinir...")
	return db.Exec(`ALTER TABLE tours DROP FOREIGN KEY fk_tour_room, DROP COLUMN room_id`).Error
}

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

	// ВРЕМЕННО: догоняет миграцию 035 на проде, где к ней нет ручного доступа
	// (mysql CLI). Убрать этот блок после того, как убедимся, что тур создаётся
	// без ошибки fk_tour_room.
	if err := applyMigration035(db); err != nil {
		log.Fatalf("migrasyon 035 hatası: %v", err)
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

	defaultTaskRepo := repository.NewDefaultTaskRepository(db)
	defaultTaskSvc := service.NewDefaultTaskService(defaultTaskRepo)

	taskRepo := repository.NewTaskRepository(db)

	taskCommentRepo := repository.NewTaskCommentRepository(db)
	taskCommentSvc := service.NewTaskCommentService(taskCommentRepo)

	tourRepo := repository.NewTourRepository(db)
	tourSvc := service.NewTourService(tourRepo, defaultTaskRepo, taskRepo)

	clientRepo := repository.NewClientRepository(db)
	clientSvc := service.NewClientService(clientRepo)

	settingRepo := repository.NewSettingRepository(db)
	settingSvc := service.NewSettingService(settingRepo)

	referansUserRepo := repository.NewReferansUserRepository(db)
	referansUserSvc := service.NewReferansUserService(referansUserRepo)

	hocaUserRepo := repository.NewHocaUserRepository(db)
	hocaUserSvc := service.NewHocaUserService(hocaUserRepo)

	// Задачи уведомляют исполнителя в Telegram, поэтому создаются после hoca и настроек.
	telegramSvc := service.NewTelegramService(settingRepo, hocaUserRepo)
	taskSvc := service.NewTaskService(taskRepo, telegramSvc)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	sessions := auth.NewStore()

	discountCategoryRepo := repository.NewDiscountCategoryRepository(db)
	discountCategorySvc := service.NewDiscountCategoryService(discountCategoryRepo)

	discountRepo := repository.NewDiscountRepository(db)
	discountSvc := service.NewDiscountService(discountRepo)

	orderRepo := repository.NewOrderRepository(db)
	orderSvc := service.NewOrderService(orderRepo, incomeRepo, discountRepo, tourRepo)

	sheetsClient, err := googlesheets.NewClient(context.Background(), googlesheets.CredentialsPath())
	if err != nil {
		log.Printf("Google Sheets qoşulmadı: %v", err)
	}

	sheetLinkRepo := repository.NewSheetLinkRepository(db)
	sheetLinkSvc := service.NewSheetLinkService(sheetLinkRepo)

	metaAdAccountRepo := repository.NewMetaAdAccountRepository(db)
	metaAdSpendRepo := repository.NewMetaAdSpendRepository(db)
	metaAdsSvc := service.NewMetaAdsService(metaAdAccountRepo, metaAdSpendRepo, expenseRepo)

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
	referansUserHandler := handler.NewReferansUserHandler(referansUserSvc)
	hocaUserHandler := handler.NewHocaUserHandler(hocaUserSvc)
	defaultTaskHandler := handler.NewDefaultTaskHandler(defaultTaskSvc)
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(userSvc, sessions)
	discountCategoryHandler := handler.NewDiscountCategoryHandler(discountCategorySvc)
	discountHandler := handler.NewDiscountHandler(discountSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	taskHandler := handler.NewTaskHandler(taskSvc)
	taskCommentHandler := handler.NewTaskCommentHandler(taskCommentSvc)
	telegramHandler := handler.NewTelegramHandler(telegramSvc)
	pageHandler := handler.NewPageHandler(accountSvc, incomeCategorySvc, incomeSvc, expenseCategorySvc, expenseSvc, tourCategorySvc, roomSvc, flightSvc, tourSvc, clientSvc, settingSvc, referansUserSvc, hocaUserSvc, defaultTaskSvc, userSvc, discountCategorySvc, discountSvc, orderSvc, taskSvc, metaAdsSvc, telegramSvc)
	sheetsImportHandler := handler.NewSheetsImportHandler(sheetsClient, sheetLinkSvc, tourSvc, clientSvc, orderSvc)
	metaAdsHandler := handler.NewMetaAdsHandler(metaAdsSvc)

	router := handler.NewRouter(authHandler, sessions, accountHandler, incomeCategoryHandler, incomeHandler, expenseCategoryHandler, expenseHandler, tourCategoryHandler, roomHandler, flightHandler, tourHandler, clientHandler, settingHandler, referansUserHandler, hocaUserHandler, defaultTaskHandler, userHandler, discountCategoryHandler, discountHandler, orderHandler, taskHandler, taskCommentHandler, telegramHandler, pageHandler, sheetsImportHandler, metaAdsHandler, tmpl)
	router.Run(":" + os.Getenv("PORT"))
}
