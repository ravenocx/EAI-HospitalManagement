package controller

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ravenocx/hospital-mgt/models"
	"github.com/ravenocx/hospital-mgt/responses"
	"github.com/ravenocx/hospital-mgt/service"
	"github.com/ravenocx/hospital-mgt/utils"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

type registerResponse struct {
	UserId      string `json:"userId"`
	Nip         int64  `json:"nip"`
	Name        string `json:"name"`
	Token       token  `json:"token"`
}

type loginResponse struct {
	UserId      string `json:"userId"`
	Nip         int64  `json:"nip"`
	Name        string `json:"name"`
	Token       token  `json:"token"`
}

type token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	var newUser models.AdminRegistrationPayload
	if err := ctx.BodyParser(&newUser); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	context := context.Background()

	userId, tokens, err := c.service.Register(context, newUser)
	if (err != responses.CustomError{}) {
		return ctx.Status(err.Status()).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	responseData := registerResponse{
		UserId:      userId,
		Nip:         newUser.Nip,
		Name:        newUser.Name,
		Token: token{
			AccessToken: tokens.Access,
			RefreshToken: tokens.Refresh,
		},
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data":    responseData,
	})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var user models.AdminCredential
	if err := ctx.BodyParser(&user); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	loginPayload := models.Credential{
		Nip:      strconv.FormatInt(user.Nip, 10),
		Password: user.Password,
	}
	context := context.Background()

	userId, name, tokens, err := c.service.Login(context, loginPayload)
	if (err != responses.CustomError{}) {
		return ctx.Status(err.Status()).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	responseData := loginResponse{
		UserId:      userId,
		Nip:         user.Nip,
		Name:        name,
		Token: token{
			AccessToken: tokens.Access,
			RefreshToken: tokens.Refresh,
		},
	}

	return ctx.JSON(fiber.Map{
		"message": "User logged in successfully",
		"data":    responseData,
	})
}

func (c *UserController) RenewTokens(ctx *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := utils.ExtractTokenMetadata(ctx)
	if err != nil {
		return responses.NewInternalServerError(err.Error())
	}

	expiresAccessToken := claims.Expires

	if now > expiresAccessToken {
		return responses.NewUnauthorizedError("your token is not expired yet")
	}

	renew := &models.Renew{}

	if err := ctx.BodyParser(renew); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	context := context.Background()

	log.Printf("token : %+v", claims.UserID)

	tokens, customError := c.service.UpdateRefreshToken(context, renew.RefreshToken, claims)
	if (customError != responses.CustomError{}) {
		log.Printf("t : %+v", customError)
		return ctx.Status(customError.Status()).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	responseData := utils.Tokens{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"tokens":  responseData,
	})
}

func (c *UserController) NurseLogin(ctx *fiber.Ctx) error {
	var user models.NurseCredential
	if err := ctx.BodyParser(&user); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	loginPayload := models.Credential{
		Nip:      strconv.FormatInt(user.Nip, 10),
		Password: user.Password,
	}

	context := context.Background()

	userId, name, tokens, err := c.service.NurseLogin(context, loginPayload)
	if (err != responses.CustomError{}) {
		return ctx.Status(err.Status()).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	responseData := loginResponse{
		UserId:      userId,
		Nip:         user.Nip,
		Name:        name,
		Token: token{
			AccessToken: tokens.Access,
			RefreshToken: tokens.Refresh,
		},
	}

	return ctx.JSON(fiber.Map{
		"message": "User logged in successfully",
		"data":    responseData,
	})
}
