package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ravenocx/hospital-mgt/models"
	"github.com/ravenocx/hospital-mgt/repositories"
	"github.com/ravenocx/hospital-mgt/responses"
	"github.com/ravenocx/hospital-mgt/utils"
)

type UserService interface {
	Register(ctx context.Context, newUser models.AdminRegistrationPayload) (string, string, responses.CustomError)
	Login(ctx context.Context, user models.Credential) (string, string, string, responses.CustomError)
	NurseLogin(ctx context.Context, creds models.Credential) (string, string, string, responses.CustomError)
	UpdateRefreshToken(ctx context.Context, refreshToken string, token *utils.TokenMetadata) (*utils.Tokens, responses.CustomError)
}

type userService struct {
	repo repositories.UserRepositories
}

func NewUserService(repo repositories.UserRepositories) UserService {
	return &userService{repo}
}

func (s *userService) Register(ctx context.Context, newUser models.AdminRegistrationPayload) (string, string, responses.CustomError) {
	validate := utils.NewValidator()

	if err := validate.Struct(&newUser); err != nil {
		return "", "", responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
	}

	existingUser, err := s.repo.GetUser(ctx, strconv.FormatInt(newUser.Nip, 10))
	if err != nil {
		if err != pgx.ErrNoRows {
			return "", "", responses.NewInternalServerError(fmt.Sprintf("failed to get existing user : %+v", err.Error()))
		}
	}

	if existingUser != nil {
		return "", "", responses.NewConflictError("user already exists")
	}

	newUser.Password = utils.GeneratePassword(newUser.Password)

	// // Use a transaction for creating the user and generating the token
	// tx, err := s.repo.BeginTx(ctx)
	// if err != nil {
	// 	return "", "", err
	// }
	// defer tx.Rollback(ctx)

	userId := uuid.New()

	userRole := repositories.CheckRoleForRegister(strconv.FormatInt(newUser.Nip, 10))

	tokens, err := utils.GenerateNewTokens(userId.String(), userRole)
	if err != nil {
		return "", "", responses.NewInternalServerError(fmt.Sprintf("failed to generate new JWT token : %+v", err.Error()))
	}

	id, err := s.repo.CreateUser(ctx, &newUser, userId, tokens.Refresh)
	if err != nil {
		return "", "", responses.NewInternalServerError(fmt.Sprintf("failed to create user : %+v", err.Error()))
	}

	// if err := tx.Commit(ctx); err != nil {
	// 	return "", "", err
	// }

	return id, tokens.Access, responses.CustomError{}
}

func (s *userService) Login(ctx context.Context, creds models.Credential) (string, string, string, responses.CustomError) {
	if strings.HasPrefix(creds.Nip, "303") {
		return "", "", "", responses.NewNotFoundError("user is not from admin (nip not starts with 615)")
	}

	validate := utils.NewValidator()

	if err := validate.Struct(&creds); err != nil {
		return "", "", "", responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
	}

	user, err := s.repo.GetUser(ctx, creds.Nip)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", "", responses.NewNotFoundError("user not found")
		}
		return "", "", "", responses.NewInternalServerError(fmt.Sprintf("failed to get user : %+v", err.Error()))
	}

	if err := utils.ComparePasswords(user.Password, creds.Password); err != nil {
		return "", "", "", responses.NewBadRequestError("wrong password, please try again!")
	}

	tokens, err := utils.GenerateNewTokens(user.ID, user.Role)
	if err != nil {
		return "", "", "", responses.NewInternalServerError(fmt.Sprintf("failed to generate new JWT token : %+v", err))
	}

	res, err := s.repo.UpdateRefreshToken(ctx, user.ID, tokens.Refresh)
	if res.RowsAffected() == 0 {
		return "", "", "", responses.NewNotFoundError(fmt.Sprintf("user not found : %+v", err.Error()))
	}

	if err != nil {
		return "", "", "", responses.NewInternalServerError(fmt.Sprintf("failed to update refresh token : %+v", err.Error()))
	}

	return user.ID, user.Name, tokens.Access, responses.CustomError{}
}

func (s *userService) NurseLogin(ctx context.Context, creds models.Credential) (string, string, string, responses.CustomError) {
	if strings.HasPrefix(creds.Nip, "615") {
		return "", "", "", responses.NewNotFoundError("user is not from nurse (nip not starts with 303)")
	}

	validate := utils.NewValidator()

	if err := validate.Struct(&creds); err != nil {
		return "", "", "", responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
	}

	userAccess, err := s.repo.GetNurseAccessByNip(ctx, creds.Nip)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", "", responses.NewNotFoundError("user not found")
		}
		return "", "", "", responses.NewInternalServerError(fmt.Sprintf("failed to get user : %+v", err.Error()))
	}
	if !userAccess.Access {
		return "", "", "", responses.NewBadRequestError("user doesn't have access")
	}

	user, err := s.repo.GetUser(ctx, creds.Nip)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", "", responses.NewNotFoundError("user not found")
		}
		return "", "", "", responses.NewInternalServerError(fmt.Sprintf("failed to get user : %+v", err.Error()))
	}


	if err := utils.ComparePasswords(user.Password, creds.Password); err != nil {
		return "", "", "", responses.NewBadRequestError("wrong password!")
	}

	tokens, err := utils.GenerateNewTokens(user.ID, user.Role)
	if err != nil {
		return "", "", "", responses.NewInternalServerError(fmt.Sprintf("failed to generate new JWT token : %+v", err))
	}

	res, err := s.repo.UpdateRefreshToken(ctx, user.ID, tokens.Refresh)
	if res.RowsAffected() == 0 {
		return "", "", "", responses.NewNotFoundError(fmt.Sprintf("user not found : %+v", err.Error()))
	}

	if err != nil {
		return "", "", "", responses.NewInternalServerError(fmt.Sprintf("failed to update refresh token : %+v", err.Error()))
	}

	return user.ID, user.Name, tokens.Access, responses.CustomError{}
}

func (s *userService) UpdateRefreshToken(ctx context.Context, refreshToken string, token *utils.TokenMetadata) (*utils.Tokens, responses.CustomError) {
	expiresRefreshToken, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
	}

	now := time.Now().Unix()

	if now < expiresRefreshToken {
		userID := token.UserID

		user, err := s.repo.GetUserById(ctx, userID.String())
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, responses.NewNotFoundError("user not found")
			}
			return nil, responses.NewInternalServerError(fmt.Sprintf("failed to get user : %+v", err.Error()))
		}

		log.Printf("test %+v", user)

		tokens, err := utils.GenerateNewTokens(user.ID, user.Role)
		if err != nil {
			return nil, responses.NewInternalServerError(fmt.Sprintf("failed to generate new JWT token : %+v", err))
		}

		res, err := s.repo.UpdateRefreshToken(ctx, user.ID, tokens.Refresh)
		if res.RowsAffected() == 0 {
			return nil, responses.NewNotFoundError(fmt.Sprintf("user not found : %+v", err.Error()))
		}
		log.Println("test")

		if err != nil {
			return nil, responses.NewInternalServerError(fmt.Sprintf("failed to update refresh token : %+v", err.Error()))
		}

		return tokens, responses.CustomError{}
	} else {
		return nil, responses.NewUnauthorizedError("unauthorized, your session was ended earlier")
	}
}
