package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ravenocx/hospital-mgt/models"
)

type MedicalRecordRepositories interface {
	GetPatient(ctx context.Context, patientIdentityNumber int64) (string, error)
	CreateRecord(ctx context.Context, patient *models.RecordRegistrationPayload, createdBy *models.CreatedByDetail) error
	GetRecord(ctx context.Context, filter models.GetRecordQueries) ([]models.GetRecordResponse, error)
}

type medicalRecordRepositories struct {
	db *pgxpool.Pool
}

func NewMedicalRecordRepo(db *pgxpool.Pool) MedicalRecordRepositories {
	return &medicalRecordRepositories{db}
}

func (r *medicalRecordRepositories) GetPatient(ctx context.Context, patientIdentityNumber int64) (string, error) {
	var identityNumber string
	query := "SELECT identity_number FROM patients WHERE identity_number = $1"

	row := r.db.QueryRow(ctx, query, patientIdentityNumber)
	err := row.Scan(&identityNumber)
	if err != nil {
		return "", err
	}

	return identityNumber, nil
}

func (r *medicalRecordRepositories) CreateRecord(ctx context.Context, record *models.RecordRegistrationPayload, createdBy *models.CreatedByDetail) error {
	statement := "INSERT INTO medical_records (identity_number, symptoms, medications, created_by_nip, created_by_name, created_by_user_id) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.Exec(ctx, statement, record.IdentityNumber, record.Symptoms, record.Medications, createdBy.Nip, createdBy.Name, createdBy.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *medicalRecordRepositories) GetRecord(ctx context.Context, filter models.GetRecordQueries) ([]models.GetRecordResponse, error) {
	var records []models.GetRecordResponse
	var createdAt time.Time

	query := "SELECT identity_number, symptoms, medications, created_by_nip, created_by_name, created_by_user_id, created_at FROM medical_records"

	query += getRecordConstructWhereQuery(filter)

	if filter.CreatedAt != "" {
		if filter.CreatedAt == "asc" {
			query += " ORDER BY created_at ASC"
		} else if filter.CreatedAt == "desc" {
			query += " ORDER BY created_at DESC"
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	query += " limit $1 offset $2"

	rows, err := r.db.Query(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		record := models.GetRecordResponse{}
		var identityNumber int64
		var birthDate time.Time
		var nipString string

		err := rows.Scan(&identityNumber, &record.Symptoms, &record.Medications, &nipString, &record.CreatedBy.Name, &record.CreatedBy.UserId, &createdAt)
		if err != nil {
			return nil, err
		}

		queryPatient := "SELECT phone_number, name, birth_date, gender, identity_card_scan_img FROM patients WHERE identity_number = $1"

		row := r.db.QueryRow(ctx, queryPatient, identityNumber)
		err = row.Scan(&record.IdentityDetail.PhoneNumber, &record.IdentityDetail.Name, &birthDate, &record.IdentityDetail.Gender, &record.IdentityDetail.IdentityCardScanImg)
		if err != nil {
			return []models.GetRecordResponse{}, err
		}

		record.CreatedBy.Nip = nipString
		if err != nil {
			return []models.GetRecordResponse{}, err
		}
		record.IdentityDetail.IdentityNumber = identityNumber
		record.IdentityDetail.BirthDate = birthDate.Format(time.RFC3339Nano)

		record.CreatedAt = createdAt.Format(time.RFC3339Nano)

		records = append(records, record)
	}

	return records, nil
}

func getRecordConstructWhereQuery(filter models.GetRecordQueries) string {
	whereSQL := []string{}

	if filter.IdentityNumber != nil {
		whereSQL = append(whereSQL, fmt.Sprintf(" identity_number = %d", *filter.IdentityNumber))
	}

	if filter.CreatedByNip != "" {
		whereSQL = append(whereSQL, " created_by_nip = '"+filter.CreatedByNip+"'")
	}

	if filter.CreatedByUserId != "" {
		whereSQL = append(whereSQL, " created_by_user_id = '"+filter.CreatedByUserId+"'")
	}

	if len(whereSQL) > 0 {
		return " WHERE " + strings.Join(whereSQL, " AND ")
	}

	return ""
}
