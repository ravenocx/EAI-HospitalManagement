package controller

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ravenocx/hospital-mgt/models"
	"github.com/ravenocx/hospital-mgt/responses"
	"github.com/ravenocx/hospital-mgt/service"
)

type PatientController struct {
	service service.PatientService
}

func NewUserController(service service.PatientService) *PatientController {
	return &PatientController{service: service}
}

type registerPatientResponse struct {
	IdentityNumber int64  `json:"identityNumber"`
	Name           string `json:"name"`
}

func (c *PatientController) RegisterPatient(ctx *fiber.Ctx) error {
	var newPatient models.PatientRegistrationPayload
	if err := ctx.BodyParser(&newPatient); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	imgFile, err := ctx.FormFile("identityCardScanImg")
	if err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	newPatient.IdentityCardScanImg = imgFile

	context := context.Background()
	custErr := c.service.RegisterPatient(context, newPatient)
	if (custErr != responses.CustomError{}) {
		return ctx.Status(custErr.Status()).JSON(fiber.Map{
			"message": custErr.Error(),
		})
	}

	responseData := registerPatientResponse{
		IdentityNumber: newPatient.IdentityNumber,
		Name:           newPatient.Name,
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data":    responseData,
	})
}

func (c *PatientController) GetPatient(ctx *fiber.Ctx) error {
	identNumberQuery, err := strconv.ParseInt(ctx.Query("identityNumber"), 10, 64)
	identNumber := &identNumberQuery

	if err != nil {
		identNumber = nil
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "5"))
	if err != nil || limit < 0 {
		limit = 5
	}
	log.Printf("limit : %+v", limit)


	offset, err := strconv.Atoi(ctx.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	name := strings.ToLower(ctx.Query("name"))
	phoneNumber := ctx.Query("phoneNumber")
	createdAt := ctx.Query("createdAt")

	patientQuery := models.GetPatientQueries{
		IdentityNumber: identNumber,
		Limit:          limit,
		Offset:         offset,
		Name:           name,
		PhoneNumber:    phoneNumber,
		CreatedAt:      createdAt,
	}

	context := context.Background()
	resp, custErr := c.service.GetPatient(context, patientQuery)
	if (custErr != responses.CustomError{}) {
		return ctx.Status(custErr.Status()).JSON(fiber.Map{
			"message": custErr.Error(),
		})
	}

	if len(resp) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
			"data":    []interface{}{},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    resp,
	})
}
