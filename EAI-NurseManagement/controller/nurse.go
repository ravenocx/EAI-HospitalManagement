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

type UserController struct {
	service service.NurseService
}

func NewUserController(service service.NurseService) *UserController {
	return &UserController{service: service}
}

type registerResponse struct {
	UserId      string `json:"userId"`
	Nip         int64  `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken,omitempty"`
}

func (c *UserController) NurseRegister(ctx *fiber.Ctx) error {
	var newNurse models.NurseRegistrationPayload
	if err := ctx.BodyParser(&newNurse); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	imgFile, err := ctx.FormFile("identityCardScanImg")
	if err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	newNurse.IdentityCardScanImg = imgFile

	context := context.Background()
	userId, custErr := c.service.NurseRegister(context, newNurse)
	if (custErr != responses.CustomError{}) {
		return ctx.Status(custErr.Status()).JSON(fiber.Map{
			"message": custErr.Error(),
		})
	}

	responseData := registerResponse{
		UserId: userId,
		Nip:    newNurse.Nip,
		Name:   newNurse.Name,
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Nurse registered successfully",
		"data":    responseData,
	})
}

func (c *UserController) NurseUpdate(ctx *fiber.Ctx) error {
	id := ctx.Params("userId")
	var updatePayload models.NurseUpdatePayload
	if err := ctx.BodyParser(&updatePayload); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	context := context.Background()
	err := c.service.UpdateNurse(context, id, updatePayload)
	if (err != responses.CustomError{}) {
		return ctx.Status(err.Status()).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":      id,
		"message": "success updated nurse",
	})
}

func (c *UserController) NurseDelete(ctx *fiber.Ctx) error {
	id := ctx.Params("userId")

	context := context.Background()
	err := c.service.DeleteNurse(context, id)

	if (err != responses.CustomError{}) {
		return ctx.Status(err.Status()).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":      id,
		"message": "success deleted nurse",
	})
}

func (c *UserController) NurseAccess(ctx *fiber.Ctx) error {
	id := ctx.Params("userId")
	var accessPayload models.NurseAccessPayload
	if err := ctx.BodyParser(&accessPayload); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	context := context.Background()
	err := c.service.AccessNurse(context, id, accessPayload)

	if (err != responses.CustomError{}) {
		return ctx.Status(err.Status()).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":      id,
		"message": "success add new access",
	})
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	userId := ctx.Query("userId")

	limit, err := strconv.Atoi(ctx.Query("limit", "5"))
	if err != nil || limit < 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(ctx.Query("limit", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	name := strings.ToLower(ctx.Query("name"))
	nip := ctx.Query("nip")
	_, err = strconv.ParseInt(nip, 10, 64)
	if err != nil {
		nip = ""
	}

	role := ctx.Query("role")
	createdAt := ctx.Query("createdAt")

	userQuery := models.GetUserQueries{
		UserId:    userId,
		Limit:     limit,
		Offset:    offset,
		Name:      name,
		Nip:       nip,
		Role:      role,
		CreatedAt: createdAt,
	}

	log.Println(userQuery.UserId)

	context := context.Background()
	resp, custErr := c.service.GetUser(context, userQuery)
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
