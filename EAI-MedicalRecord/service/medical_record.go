package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/ravenocx/hospital-mgt/models"
	"github.com/ravenocx/hospital-mgt/repositories"
	"github.com/ravenocx/hospital-mgt/responses"
	"github.com/ravenocx/hospital-mgt/utils"
)

type MedicalRecordService interface {
	RegisterRecord(ctx context.Context, newRecord models.RecordRegistrationPayload, createdByDetail models.CreatedByDetail, jwtToken string) responses.CustomError
	GetRecord(ctx context.Context, GetRecordQueries models.GetRecordQueries) ([]models.GetRecordResponse, responses.CustomError)
	GetNurseDetail(nurseId string, jwtToken string) ([]models.Nurse, responses.CustomError)
}

type medicalRecordService struct {
	repo repositories.MedicalRecordRepositories
}

func NewMedicalServiceService(repo repositories.MedicalRecordRepositories) MedicalRecordService {
	return &medicalRecordService{repo}
}

func (s *medicalRecordService) RegisterRecord(ctx context.Context, newRecord models.RecordRegistrationPayload, createdByDetail models.CreatedByDetail, jwtToken string) responses.CustomError {
	validate := utils.NewValidator()

	if err := validate.Struct(&newRecord); err != nil {
		return responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
	}

	existingPatient, err := GetPatient(newRecord.IdentityNumber, jwtToken) // TODO : get patient should consume endpoint get patientn
	if err != nil {
		if err.Error() == "patient with identityNumber is not exist" {
			return responses.NewNotFoundError("patient with identity_number is not exist")
		}

		log.Println(err.Error())
		return responses.NewInternalServerError(err.Error())

	}

	if existingPatient == nil {
		return responses.NewNotFoundError("patient with identity_number is not exist")
	}

	err = s.repo.CreateRecord(ctx, &newRecord, &createdByDetail)
	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to create new medical record : %+v", err.Error()))
	}

	return responses.CustomError{}
}

func (s *medicalRecordService) GetRecord(ctx context.Context, GetRecordQueries models.GetRecordQueries) ([]models.GetRecordResponse, responses.CustomError) {

	validate := utils.NewValidator()

	if GetRecordQueries.IdentityNumber != nil {
		if err := validate.Struct(&GetRecordQueries); err != nil {
			return nil, responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
		}
	}

	patients, err := s.repo.GetRecord(ctx, GetRecordQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.GetRecordResponse{}, responses.CustomError{}
		}
		return nil, responses.NewInternalServerError(fmt.Sprintf("failed to get medical record : %+v", err.Error()))
	}

	return patients, responses.CustomError{}
}

func GetPatient(identityNumber int64, jwtToken string) ([]models.Patient, error) {
	medicalUserUrl := "http://localhost:5000/v1/medical/patient"
	params := url.Values{}
	params.Add("identityNumber", strconv.FormatInt(identityNumber, 10))

	reqURL := fmt.Sprintf("%s?%s", medicalUserUrl, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, responses.NewNotFoundError("patient with identityNumber is not exist")
		}
		return nil, responses.NewInternalServerError("failed to consume get patient endpoint")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var patient models.PatientResponse
	err = json.Unmarshal(body, &patient)
	if err != nil {
		log.Printf("errJson : %+v", err)
		return nil, err
	}

	// if patient.Message == "success" {
	// 	for _, patient := range patient.Data {
	// 		fmt.Printf("User ID: %s, NIP: %s, Name: %s", )
	// 	}

	// 	// If you only want the first (newest) patient's details:
	// 	if len(patient.Data) > 0 {
	// 		firstpatient := patient.Data[0]
	// 		fmt.Printf("Newest patient - User ID: %s, NIP: %s, Name: %s\n", firstpatient.UserID, firstpatient.NIP, firstpatient.Name)
	// 	}
	// } else {
	// 	log.Println("Failed to fetch patient details")
	// }

	return patient.Data, nil
}

func (s *medicalRecordService) GetNurseDetail(nurseId string, jwtToken string) ([]models.Nurse, responses.CustomError) {
	nurseUrl := "http://localhost:4000/v1/user"
	params := url.Values{}
	params.Add("userId", nurseId)

	reqURL := fmt.Sprintf("%s?%s", nurseUrl, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, responses.NewInternalServerError(err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, responses.NewInternalServerError(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, responses.NewNotFoundError("nurse with nurse_id is not exist")
		}
		log.Println(resp.StatusCode)
		return nil, responses.NewInternalServerError("failed to consume get user endpoint")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, responses.NewInternalServerError(err.Error())
	}

	var nurse models.NurseResponse
	err = json.Unmarshal(body, &nurse)
	if err != nil {
		return nil, responses.NewInternalServerError(err.Error())
	}

	// if nurse.Message == "success" {
	// 	for _, nurse := range nurse.Data {
	// 		fmt.Printf("User ID: %s, NIP: %s, Name: %s\n", nurse.UserID, nurse.NIP, nurse.Name)
	// 	}

	// 	// If you only want the first (newest) nurse's details:
	// 	if len(nurse.Data) > 0 {
	// 		firstNurse := nurse.Data[0]
	// 		fmt.Printf("Newest Nurse - User ID: %s, NIP: %s, Name: %s\n", firstNurse.UserID, firstNurse.NIP, firstNurse.Name)
	// 	}
	// } else {
	// 	log.Println("Failed to fetch nurse details")
	// }

	return nurse.Data, responses.CustomError{}
}
