package repositories

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ravenocx/hospital-mgt/models"
)

type NurseRepositories interface {
	GetUser(ctx context.Context, nip string) (*models.User, error)
	CreateNurseUser(ctx context.Context, user *models.NurseRegistrationPayload) (string, error)
	GetUserNipById(ctx context.Context, id string) (*models.User, error)
	UpdateNurse(ctx context.Context, nurseId string, updatePayload models.NurseUpdatePayload) (pgconn.CommandTag, error)
	UpdateAccessNurse(ctx context.Context, nurseId string, passwordHash string) (pgconn.CommandTag, error)
	DeleteNurse(ctx context.Context, userId string) (pgconn.CommandTag, error)
	GetUsers(ctx context.Context, filter models.GetUserQueries) ([]models.GetUserResponse, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type nurseRepositories struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) NurseRepositories {
	return &nurseRepositories{db}
}

func (r *nurseRepositories) GetUser(ctx context.Context, nip string) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, access, role FROM users WHERE nip = $1"

	row := r.db.QueryRow(ctx, query, nip)
	err := row.Scan(&user.ID, &user.Name, &user.Access, &user.Role)
	if err != nil {
		return nil, err
	}
	user.Nip = nip

	return &user, nil
}

func (r *nurseRepositories) CreateNurseUser(ctx context.Context, user *models.NurseRegistrationPayload) (string, error) {
	var id string
	role := CheckRoleForRegister(strconv.FormatInt(user.Nip, 10))
	statement := "INSERT INTO users (name, nip, role, access, identity_card_scan_img) VALUES ($1, $2, $3, false, $4) RETURNING id"

	row := r.db.QueryRow(ctx, statement, user.Name, strconv.FormatInt(user.Nip, 10), role, user.IdentityCardScanImgString)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *nurseRepositories) GetUserNipById(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := "SELECT nip FROM users WHERE id = $1"

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&user.Nip)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *nurseRepositories) UpdateNurse(ctx context.Context, nurseId string, updatePayload models.NurseUpdatePayload) (pgconn.CommandTag, error) {
	statement := "UPDATE users SET nip = $1, name = $2, updated_at = $3 WHERE id = $4"

	res, err := r.db.Exec(ctx, statement, strconv.FormatInt(updatePayload.Nip, 10), updatePayload.Name, time.Now(), nurseId)

	return res, err
}

func (r *nurseRepositories) DeleteNurse(ctx context.Context, userId string) (pgconn.CommandTag, error) {
	statement := "DELETE FROM users WHERE id = $1"

	res, err := r.db.Exec(ctx, statement, userId)
	return res, err
}

func (r *nurseRepositories) UpdateAccessNurse(ctx context.Context, nurseId string, passwordHash string) (pgconn.CommandTag, error) {
	statement := "UPDATE users SET password = $1, access = true WHERE id = $2"

	res, err := r.db.Exec(ctx, statement, passwordHash, nurseId)

	return res, err
}

func (r *nurseRepositories) GetUsers(ctx context.Context, filter models.GetUserQueries) ([]models.GetUserResponse, error) {
	var users []models.GetUserResponse
	var createdAt time.Time
	query := "SELECT id, nip, name, created_at FROM users"

	query += getUserConstructWhereQuery(filter)

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
	log.Printf("Get users query : %+v", query)

	rows, err := r.db.Query(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	var nip string
	for rows.Next() {
		user := models.GetUserResponse{}
		err := rows.Scan(&user.UserId, &nip, &user.Name, &createdAt)
		if err != nil {
			return nil, err
		}
		user.Nip, err = strconv.ParseInt(nip, 10, 64)
		if err != nil {
			return nil, err
		}
		user.CreatedAt = createdAt.Format(time.RFC3339Nano)
		users = append(users, user)
	}

	return users, nil
}

func (r *nurseRepositories) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

func CheckRoleForRegister(nip string) string {
	if strings.HasPrefix(nip, "615") {
		return "admin"
	} else if strings.HasPrefix(nip, "303") {
		return "nurse"
	}
	return ""
}

func getUserConstructWhereQuery(filter models.GetUserQueries) string {
	whereSQL := []string{}

	if filter.UserId != "" {
		whereSQL = append(whereSQL, " id = '"+filter.UserId+"'")
	}

	if filter.Name != "" {
		whereSQL = append(whereSQL, " name ILIKE '%"+filter.Name+"%'")
	}

	if filter.Nip != "" {
		whereSQL = append(whereSQL, " nip ILIKE '"+filter.Nip+"%'")
	}

	if filter.Role == "admin" {
		whereSQL = append(whereSQL, " role = '"+filter.Role+"'")
	}

	if filter.Role == "nurse" {
		whereSQL = append(whereSQL, " role = '"+filter.Role+"'")
	}

	if len(whereSQL) > 0 {
		return " WHERE " + strings.Join(whereSQL, " AND ")
	}

	return ""
}
