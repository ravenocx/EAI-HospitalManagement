package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ravenocx/hospital-mgt/models"
	"github.com/ravenocx/hospital-mgt/repositories"
	"github.com/ravenocx/hospital-mgt/responses"
	"github.com/ravenocx/hospital-mgt/utils"
)

type NurseService interface {
	NurseRegister(ctx context.Context, newUser models.NurseRegistrationPayload) (string, responses.CustomError)
	UpdateNurse(ctx context.Context, nurseId string, updatePayload models.NurseUpdatePayload) responses.CustomError
	DeleteNurse(ctx context.Context, nurseId string) responses.CustomError
	AccessNurse(ctx context.Context, nurseId string, password models.NurseAccessPayload) responses.CustomError
	GetUser(ctx context.Context, GetUserQueries models.GetUserQueries) ([]models.GetUserResponse, responses.CustomError)
}

type nurseService struct {
	repo repositories.NurseRepositories
}

func NewNurseService(repo repositories.NurseRepositories) NurseService {
	return &nurseService{repo}
}

func (s *nurseService) NurseRegister(ctx context.Context, newUser models.NurseRegistrationPayload) (string, responses.CustomError) {
	validate := utils.NewValidator()

	if err := validate.Struct(&newUser); err != nil {
		return "", responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
	}

	existingUser, err := s.repo.GetUser(ctx, strconv.FormatInt(newUser.Nip, 10))
	if err != nil {
		if err != pgx.ErrNoRows {
			return "", responses.NewInternalServerError(fmt.Sprintf("failed to get existing user : %+v", err.Error()))
		}
	}

	if existingUser != nil {
		return "", responses.NewConflictError("user already exists")
	}

	urlImage, err := utils.UploadImage(newUser.IdentityCardScanImg)
	if err != nil {
		return "", responses.NewInternalServerError(fmt.Sprintf("failed to upload image : %+v", err.Error()))
	}

	log.Printf("urlImage : %+v", urlImage)

	newUser.IdentityCardScanImgString = urlImage

	id, err := s.repo.CreateNurseUser(ctx, &newUser)
	if err != nil {
		return "", responses.NewInternalServerError(fmt.Sprintf("failed to create user : %+v", err.Error()))
	}

	return id, responses.CustomError{}
}

func (s *nurseService) UpdateNurse(ctx context.Context, nurseId string, updatePayload models.NurseUpdatePayload) responses.CustomError {

	validate := utils.NewValidator()

	if err := validate.Struct(&updatePayload); err != nil {
		return responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
	}

	if _, err := uuid.Parse(nurseId); err != nil {
		return responses.NewNotFoundError("nurse not found or userId is not in valid format")
	}

	existingUser, err := s.repo.GetUser(ctx, strconv.FormatInt(updatePayload.Nip, 10))
	if err != nil {
		if err != pgx.ErrNoRows {
			return responses.NewInternalServerError(fmt.Sprintf("failed to get existing user : %+v", err.Error()))
		}
	}

	if existingUser != nil {
		return responses.NewConflictError("conflict, nip already used")
	}

	user, err := s.repo.GetUserNipById(ctx, nurseId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return responses.NewNotFoundError(fmt.Sprintf("user not found : %+v", err.Error()))
		}
		return responses.NewInternalServerError(fmt.Sprintf("failed to get nurse : %+v", err.Error()))
	}

	if !strings.HasPrefix(user.Nip, "303") {
		return responses.NewNotFoundError("user is not a nurse (nip not starts with 303)")
	}

	res, err := s.repo.UpdateNurse(ctx, nurseId, updatePayload)

	if res.RowsAffected() == 0 {
		return responses.NewNotFoundError(fmt.Sprintf("nurse not found : %+v", err.Error()))
	}

	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to update nurse : %+v", err.Error()))
	}

	return responses.CustomError{}
}

func (s *nurseService) DeleteNurse(ctx context.Context, nurseId string) responses.CustomError {
	if _, err := uuid.Parse(nurseId); err != nil {
		return responses.NewNotFoundError("nurse not found or userId is not in valid format")
	}

	user, err := s.repo.GetUserNipById(ctx, nurseId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return responses.NewNotFoundError(fmt.Sprintf("user not found : %+v", err.Error()))
		}
		return responses.NewInternalServerError(fmt.Sprintf("failed to get nurse : %+v", err.Error()))
	}

	if !strings.HasPrefix(user.Nip, "303") {
		return responses.NewNotFoundError("user is not a nurse (nip not starts with 303)")
	}

	res, err := s.repo.DeleteNurse(ctx, nurseId)

	if res.RowsAffected() == 0 {
		return responses.NewNotFoundError(fmt.Sprintf("nurse not found : %+v", err.Error()))
	}

	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to delete nurse : %+v", err.Error()))
	}

	return responses.CustomError{}
}

func (s *nurseService) AccessNurse(ctx context.Context, nurseId string, accessPayload models.NurseAccessPayload) responses.CustomError {
	validate := utils.NewValidator()

	if err := validate.Struct(&accessPayload); err != nil {
		return responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
	}

	if _, err := uuid.Parse(nurseId); err != nil {
		return responses.NewNotFoundError("nurse not found or userId is not in valid format")
	}

	user, err := s.repo.GetUserNipById(ctx, nurseId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return responses.NewNotFoundError(fmt.Sprintf("user not found : %+v", err.Error()))
		}
		return responses.NewInternalServerError(fmt.Sprintf("failed to get nurse : %+v", err.Error()))
	}

	if !strings.HasPrefix(user.Nip, "303") {
		return responses.NewNotFoundError("user is not a nurse (nip not starts with 303)")
	}

	hashedPassword := utils.GeneratePassword(accessPayload.Password)

	res, err := s.repo.UpdateAccessNurse(ctx, nurseId, hashedPassword)

	if res.RowsAffected() == 0 {
		return responses.NewNotFoundError(fmt.Sprintf("nurse not found : %+v", err.Error()))
	}

	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to update nurse : %+v", err.Error()))
	}

	return responses.CustomError{}
}

func (s *nurseService) GetUser(ctx context.Context, GetUserQueries models.GetUserQueries) ([]models.GetUserResponse, responses.CustomError) {

	validate := utils.NewValidator()

	if GetUserQueries.UserId != "" {
		if err := validate.Struct(&GetUserQueries); err != nil {
			return nil, responses.NewBadRequestError(fmt.Sprintf("query params doesn't meet requirement : %+v", err.Error()))
		}
	}

	users, err := s.repo.GetUsers(ctx, GetUserQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, responses.NewNotFoundError(fmt.Sprintf("user not found : %+v", err.Error()))
		}
		return nil, responses.NewInternalServerError(fmt.Sprintf("failed to get nurse : %+v", err.Error()))

	}

	return users, responses.CustomError{}
}
