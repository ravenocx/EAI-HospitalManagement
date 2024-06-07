package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ravenocx/hospital-mgt/models"
	"github.com/ravenocx/hospital-mgt/repositories"
	"github.com/ravenocx/hospital-mgt/responses"
	"github.com/ravenocx/hospital-mgt/utils"
)

type NurseService interface {
	NurseRegister(ctx context.Context, newUser models.NurseRegistrationPayload) (string, responses.CustomError)
	UpdateNurse(ctx context.Context, nurseId string, updatePayload models.NurseUpdatePayload) responses.CustomError
	DeleteNurse(ctx context.Context, nurseId string) responses.CustomError
	AccessNurse(ctx context.Context, nurseId string, password models.NurseAccessPayload) (*models.User, responses.CustomError)
	GetUser(ctx context.Context, GetUserQueries models.GetUserQueries) ([]models.GetUserResponse, responses.CustomError)
	PublishToRabbitmq(nurse *models.User) responses.CustomError
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

func (s *nurseService) AccessNurse(ctx context.Context, nurseId string, accessPayload models.NurseAccessPayload) (*models.User, responses.CustomError) {
	validate := utils.NewValidator()

	if err := validate.Struct(&accessPayload); err != nil {
		return nil, responses.NewBadRequestError(fmt.Sprintf("payload request doesn't meet requirement : %+v", err.Error()))
	}

	if _, err := uuid.Parse(nurseId); err != nil {
		return nil, responses.NewNotFoundError("nurse not found or userId is not in valid format")
	}

	user, err := s.repo.GetUserNipById(ctx, nurseId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, responses.NewNotFoundError(fmt.Sprintf("user not found : %+v", err.Error()))
		}
		return nil, responses.NewInternalServerError(fmt.Sprintf("failed to get nurse : %+v", err.Error()))
	}

	if !strings.HasPrefix(user.Nip, "303") {
		return nil, responses.NewNotFoundError("user is not a nurse (nip not starts with 303)")
	}

	hashedPassword := utils.GeneratePassword(accessPayload.Password)

	res, err := s.repo.UpdateAccessNurse(ctx, nurseId, hashedPassword)

	if res.RowsAffected() == 0 {
		return nil, responses.NewNotFoundError(fmt.Sprintf("nurse not found : %+v", err.Error()))
	}

	if err != nil {
		return nil, responses.NewInternalServerError(fmt.Sprintf("failed to update nurse : %+v", err.Error()))
	}

	return user, responses.CustomError{}
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

func (s *nurseService) PublishToRabbitmq(nurse *models.User) responses.CustomError {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to connect to rabbitmq service : %+v", err.Error()))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to create channel : %+v", err.Error()))
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"nurse_access", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil, 
	)

	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to declare rabbitmq queue : %+v", err.Error()))
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"dlx_exchange",
		false,
		amqp.Table{
			"x-message-ttl":          int32(60000), // TTL in milliseconds
			"x-dead-letter-exchange": "dlx_exchange",
		},
	)

	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to bind queue : %+v", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(nurse)

	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to encode the nurse : %+v", err.Error()))
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})

	if err != nil {
		return responses.NewInternalServerError(fmt.Sprintf("failed to publish data to rabbitmq : %+v", err.Error()))
	}

	log.Printf(" [x] Sent to rabbitmq :  %s\n", body)
	return responses.CustomError{}
}
