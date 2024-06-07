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

	PatientRoute(mainRoute, s.dbPool)
}

func PatientRoute(r fiber.Router, db *pgxpool.Pool) {
	c := controller.NewUserController(service.NewUserService(repositories.NewPatientRepo(db)))

	medicalRoute := r.Group("/medical")

	medicalRoute.Post("/patient", middleware.JWTProtected(), middleware.UserAuth(), c.RegisterPatient)
	medicalRoute.Get("/patient", middleware.JWTProtected(), middleware.UserAuth(), c.GetPatient)
}
