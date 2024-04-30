package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/kbc0/DynamicStockManager/handler"
	"github.com/kbc0/DynamicStockManager/middleware"
	user "github.com/kbc0/DynamicStockManager/repository/user"
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
	userRepo := user.NewUserRepository(srv.DB) // Assuming NewUserRepository is adjusted to accept *mongo.Database
	userHandler := handler.NewUserHandler(userRepo)
	
	// User creation endpoint
	srv.App.Post("/api/v1/register", userHandler.RegisterUser)
	srv.App.Post("/api/v1/login", userHandler.LoginUser)
}

