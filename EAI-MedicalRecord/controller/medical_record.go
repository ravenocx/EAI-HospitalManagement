package controller

import (
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ravenocx/hospital-mgt/middleware"
	"github.com/ravenocx/hospital-mgt/models"
	"github.com/ravenocx/hospital-mgt/responses"
	"github.com/ravenocx/hospital-mgt/service"
	"github.com/ravenocx/hospital-mgt/utils"
)

type MedicalRecordController struct {
	service service.MedicalRecordService
}

func NewUserController(service service.MedicalRecordService) *MedicalRecordController {
	return &MedicalRecordController{service: service}
}

func (c *MedicalRecordController) RegisterRecord(ctx *fiber.Ctx) error {
	var newRecord models.RecordRegistrationPayload
	if err := ctx.BodyParser(&newRecord); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	claims, err := utils.ExtractTokenMetadata(ctx)
	if err != nil {
		log.Println(err)
		return middleware.UnauthorizedResponse(ctx, "token not found")
	}

	userID := claims.UserID

	jwtToken := utils.ExtractToken(ctx)

	nurse, custErr := c.service.GetNurseDetail(userID.String(), jwtToken)
	if (custErr != responses.CustomError{}) {
		return ctx.Status(custErr.Status()).JSON(fiber.Map{
			"message": custErr.Error(),
		})
	}

	// TODO : get user should consume endpoint get user
	createdByDetail := models.CreatedByDetail{
		UserId: userID.String(),
		Nip:    strconv.FormatInt(nurse[0].NIP, 10) ,
		Name:   nurse[0].Name,
	}

	context := context.Background()
	custErr = c.service.RegisterRecord(context, newRecord, createdByDetail, jwtToken)
	if (custErr != responses.CustomError{}) {
		return ctx.Status(custErr.Status()).JSON(fiber.Map{
			"message": custErr.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Medical record added successfully",
		"data":    newRecord,
	})
}

func (c *MedicalRecordController) GetRecord(ctx *fiber.Ctx) error {
	identNumberQuery, err := strconv.ParseInt(ctx.Query("identityNumber"), 10, 64)
	identNumber := &identNumberQuery
	if err != nil {
		identNumber = nil
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "5"))
	if err != nil || limit < 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(ctx.Query("limit", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	nip := ctx.Query("nip")
	_, err = strconv.ParseInt(nip, 10, 64)
	if err != nil {
		nip = ""
	}

	userId := ctx.Query("userId")
	createdAt := ctx.Query("createdAt")

	recordQuery := models.GetRecordQueries{
		IdentityNumber:  identNumber,
		Limit:           limit,
		Offset:          offset,
		CreatedByNip:    nip,
		CreatedByUserId: userId,
		CreatedAt:       createdAt,
	}

	context := context.Background()
	resp, custErr := c.service.GetRecord(context, recordQuery)
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
