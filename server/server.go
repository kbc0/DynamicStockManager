package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	userHandler "github.com/kbc0/DynamicStockManager/handler/user"
	formHandler "github.com/kbc0/DynamicStockManager/handler/form"
	"github.com/kbc0/DynamicStockManager/middleware"
	userRepo "github.com/kbc0/DynamicStockManager/repository/user"
	formRepo "github.com/kbc0/DynamicStockManager/repository/form"
)

type Server struct {
	App    *fiber.App
	DB     *mongo.Database // Update this to use mongo.Database instead of gorm.DB
	logger *zerolog.Logger
}

func NewServer(db *mongo.Database, logger *zerolog.Logger) *Server {
	logger.Info().Msg("Server is created")
	app := fiber.New()

	middleware.RegisterMiddleware(app) // Ensure this is implemented to register any necessary middleware

	srv := &Server{
		App:    app,
		DB:     db,
		logger: logger,
	}

	srv.registerRoutes() // Register user routes

	return srv
}

func (srv *Server) registerRoutes() {
	// User related routes
	userRepo := userRepo.NewUserRepository(srv.DB) // Assuming NewUserRepository is adjusted to accept *mongo.Database
	userHandler := userHandler.NewUserHandler(userRepo)
	
	// User creation endpoint
	srv.App.Post("/api/v1/register", userHandler.RegisterUser)
	srv.App.Post("/api/v1/login", userHandler.LoginUser)
	srv.App.Get("/api/v1/account", userHandler.GetAccount)

	// Form related routes
	formRepo := formRepo.NewFormRepository(srv.DB)
	formHandler := formHandler.NewFormHandler(formRepo)

	// Form endpoints
	srv.App.Post("/api/v1/form/create", formHandler.CreateFormHandler)
	srv.App.Get("/api/v1/form", formHandler.GetFormsHandler)
	srv.App.Get("/api/v1/form/:_id", formHandler.GetFormHandler)
	srv.App.Put("/api/v1/form/:_id", formHandler.UpdateFormHandler)
	srv.App.Delete("/api/v1/form/:_id", formHandler.DeleteFormHandler)
}

