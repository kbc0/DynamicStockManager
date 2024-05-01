package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	userHandler "github.com/kbc0/DynamicStockManager/handler/user"
	formHandler "github.com/kbc0/DynamicStockManager/handler/form"
	fieldHandler "github.com/kbc0/DynamicStockManager/handler/field"
	stockHandler "github.com/kbc0/DynamicStockManager/handler/stock" // Import the stock handler
	"github.com/kbc0/DynamicStockManager/middleware"
	userRepo "github.com/kbc0/DynamicStockManager/repository/user"
	formRepo "github.com/kbc0/DynamicStockManager/repository/form"
	fieldRepo "github.com/kbc0/DynamicStockManager/repository/field"
	stockRepo "github.com/kbc0/DynamicStockManager/repository/stock" // Import the stock repository
)

type Server struct {
	App    *fiber.App
	DB     *mongo.Database
	logger *zerolog.Logger
}

func NewServer(db *mongo.Database, logger *zerolog.Logger) *Server {
	logger.Info().Msg("Server is created")
	app := fiber.New()

	middleware.RegisterMiddleware(app) // Assume middleware setup is already in place

	srv := &Server{
		App:    app,
		DB:     db,
		logger: logger,
	}

	srv.registerRoutes()

	return srv
}

func (srv *Server) registerRoutes() {
	// User related routes setup
	userRepo := userRepo.NewUserRepository(srv.DB)
	userHandler := userHandler.NewUserHandler(userRepo)
	srv.App.Post("/api/v1/register", userHandler.RegisterUser)
	srv.App.Post("/api/v1/login", userHandler.LoginUser)
	srv.App.Get("/api/v1/account", userHandler.GetAccount)

	// Field related routes setup
	fieldRepo := fieldRepo.NewFieldRepository(srv.DB)
	fieldHandler := fieldHandler.NewFieldHandler(fieldRepo)
	srv.App.Post("/api/v1/form/:_id/field", fieldHandler.AddFieldToForm)
	srv.App.Get("/api/v1/form/:_id/field", fieldHandler.GetAllFields)
	srv.App.Get("/api/v1/form/:_id/field/:field_id", fieldHandler.GetField)
	srv.App.Delete("/api/v1/form/:_id/field/:field_id", fieldHandler.DeleteField)
	srv.App.Put("/api/v1/form/:_id/field/:field_id", fieldHandler.UpdateField)

	// Stock related routes setup
	stockRepo := stockRepo.NewStockRepository(srv.DB)
	stockHandler := stockHandler.NewStockHandler(stockRepo, fieldRepo)
	srv.App.Post("/api/v1/form/:_id/stock", stockHandler.AddStock)
	srv.App.Get("/api/v1/form/:_id/stock", stockHandler.GetAllStocks)
	srv.App.Get("/api/v1/form/:_id/stock/:stock_id", stockHandler.GetStock)
	srv.App.Put("/api/v1/form/:_id/stock/:stock_id", stockHandler.UpdateStock)
	srv.App.Delete("/api/v1/form/:_id/stock/:stock_id", stockHandler.DeleteStock)

	// Form related routes setup
	formRepo := formRepo.NewFormRepository(srv.DB)
	formHandler := formHandler.NewFormHandler(formRepo,fieldRepo, stockRepo)
	srv.App.Post("/api/v1/form/create", formHandler.CreateFormHandler)
	srv.App.Get("/api/v1/form", formHandler.GetFormsHandler)
	srv.App.Get("/api/v1/form/:_id", formHandler.GetFormHandler)
	srv.App.Put("/api/v1/form/:_id", formHandler.UpdateFormHandler)
	srv.App.Delete("/api/v1/form/:_id", formHandler.DeleteFormHandler)

}
