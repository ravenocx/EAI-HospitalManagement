package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ravenocx/hospital-mgt/models"
	"github.com/ravenocx/hospital-mgt/repositories"
	"github.com/ravenocx/hospital-mgt/responses"
	"github.com/ravenocx/hospital-mgt/utils"
)

type PatientService interface {
	RegisterPatient(ctx context.Context, newPatient models.PatientRegistrationPayload) responses.CustomError
	GetPatient(ctx context.Context, GetPatientQueries models.GetPatientQueries) ([]models.GetPatientResponse, responses.CustomError)
}

type patientService struct {
	repo repositories.PatientRepositories
}

func NewUserService(repo repositories.PatientRepositories) PatientService {
	return &patientService{repo}
}

func (s *patientService) RegisterPatient(ctx context.Context, newPatient models.PatientRegistrationPayload) responses.CustomError {
	validate := utils.NewValidator()

	if err := validate.Struct(&newPatient); err != nil {
		return responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
	}

	existingPatient, err := s.repo.GetPatient(ctx, newPatient.IdentityNumber)
	if err != nil {
		if err != pgx.ErrNoRows {
			return responses.NewInternalServerError(fmt.Sprintf("failed to check existing user : %+v", err.Error()))
		}
	}

	if existingPatient != "" {
		return responses.NewConflictError("patient with identity number provided is already exists")
	}

	urlImage, err := utils.UploadImage(newPatient.IdentityCardScanImg)
	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to upload image : %+v", err.Error()))
	}

	newPatient.IdentityCardScanImgString = urlImage

	err = s.repo.CreatePatient(ctx, &newPatient)
	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to create patient : %+v", err.Error()))
	}

	return responses.CustomError{}
}

func (s *patientService) GetPatient(ctx context.Context, GetPatientQueries models.GetPatientQueries) ([]models.GetPatientResponse, responses.CustomError) {

	validate := utils.NewValidator()

	if GetPatientQueries.IdentityNumber != nil{
		if err := validate.Struct(&GetPatientQueries); err != nil {
			return nil, responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
		}
	}
	patients, err := s.repo.GetPatients(ctx, GetPatientQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.GetPatientResponse{}, responses.CustomError{}
		}
		return nil, responses.NewInternalServerError(fmt.Sprintf("failed to get patients : %+v", err.Error()))
	}

	return patients, responses.CustomError{}
}
