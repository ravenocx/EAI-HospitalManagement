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
	mainRoute := s.app.Group("/v1")

	MedicalRoute(mainRoute, s.dbPool)
}

func MedicalRoute(r fiber.Router, db *pgxpool.Pool) {
	c := controller.NewUserController(service.NewMedicalServiceService(repositories.NewMedicalRecordRepo(db)))

	medicalRoute := r.Group("/medical")

	medicalRoute.Post("/record", middleware.JWTProtected(), middleware.UserAuth(), c.RegisterRecord)
	medicalRoute.Get("/record", middleware.JWTProtected(), middleware.UserAuth(), c.GetRecord)
}
