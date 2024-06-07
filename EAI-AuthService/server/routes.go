package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ravenocx/hospital-mgt/controller"
	"github.com/ravenocx/hospital-mgt/middleware"
	"github.com/ravenocx/hospital-mgt/repositories"
	"github.com/ravenocx/hospital-mgt/service"
)

func (s *Server) RegisterRoute() {
	mainRoute := s.app.Group("/v1/user")

	UserRoute(mainRoute, s.dbPool)
}

func UserRoute(r fiber.Router, db *pgxpool.Pool) {
	c := controller.NewUserController(service.NewUserService(repositories.NewUserRepo(db)))

	r.Post("/nurse/login", c.NurseLogin)
	r.Post("/token/renew", middleware.JWTProtected(), c.RenewTokens)

	adminRoute := r.Group("/admin")

	adminRoute.Post("/register", c.Register)
	adminRoute.Post("/login", c.Login)


}
