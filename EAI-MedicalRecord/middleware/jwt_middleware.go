package middleware

import (
	"log"
	"os"
	"time"

	jwtMiddleware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/ravenocx/hospital-mgt/responses"
	"github.com/ravenocx/hospital-mgt/utils"
)

func JWTProtected() func(*fiber.Ctx) error {
	config := jwtMiddleware.Config{
		SigningKey:   jwtMiddleware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET_KEY"))},
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func UserAuth() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		claims, err := utils.ExtractTokenMetadata(c)
		if err != nil {
			log.Println(err)
			return UnauthorizedResponse(c, "token not found")
		}

		expires := claims.Expires
		now := time.Now().Unix()

		if now > expires {
			return UnauthorizedResponse(c, "token expired")
		}

		return c.Next()
	}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return UnauthorizedResponse(c, err.Error())
}

func UnauthorizedResponse(ctx *fiber.Ctx, message string) error {
	err := responses.NewUnauthorizedError(message)

	return ctx.Status(err.Status()).JSON(fiber.Map{
		"message": err.Error(),
	})
}
