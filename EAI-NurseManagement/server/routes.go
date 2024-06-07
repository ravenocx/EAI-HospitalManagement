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

	NurseRoute(mainRoute, s.dbPool)
}

func NurseRoute(r fiber.Router, db *pgxpool.Pool) {
	c := controller.NewUserController(service.NewNurseService(repositories.NewUserRepo(db)))

	r.Get("/",  middleware.JWTProtected(), middleware.AdminAuth(), c.GetUser)

	nurseRoute := r.Group("/nurse")

	nurseRoute.Post("/register", middleware.JWTProtected(), middleware.AdminAuth(), c.NurseRegister)
	nurseRoute.Put("/:userId", middleware.JWTProtected(), middleware.AdminAuth(), c.NurseUpdate)
	nurseRoute.Delete("/:userId", middleware.JWTProtected(), middleware.AdminAuth(), c.NurseDelete)
	nurseRoute.Post("/:userId/access", middleware.JWTProtected(), middleware.AdminAuth(), c.NurseAccess)

}
