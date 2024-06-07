package server

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ravenocx/hospital-mgt/config"
	"github.com/ravenocx/hospital-mgt/middleware"
)

type Server struct {
	dbPool *pgxpool.Pool
	app    *fiber.App
}

func NewServer(db *pgxpool.Pool, config config.Config) *Server {
	fiberConfig := fiber.Config{
		ReadTimeout: time.Duration(config.ServerReadTimeout) * time.Second,
	}

	app := fiber.New(fiberConfig)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	middleware.FiberMiddleware(app)

	return &Server{
		dbPool: db,
		app : app,
	}
}

func (s *Server) StarApp(config config.Config) {
	if err := s.app.Listen(config.ServerHost + ":" + config.ServerPort); err != nil {
		log.Fatalf("Oops... Server is not running! Reason: %v", err)
	}
}
