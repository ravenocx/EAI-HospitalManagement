package repositories

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ravenocx/hospital-mgt/models"
)

type UserRepositories interface {
	GetUser(ctx context.Context, nip string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.AdminRegistrationPayload, userId uuid.UUID, refreshToken string) (string, error)
	CreateUserTx(ctx context.Context, tx pgx.Tx, user *models.AdminRegistrationPayload, hashPassword string) (string, error)
	UpdateRefreshToken(ctx context.Context, userId string, refreshToken string) (pgconn.CommandTag, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetNurseAccessByNip(ctx context.Context, nip string) (*models.User, error)
}

type userRepositories struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepositories {
	return &userRepositories{db}
}

func (r *userRepositories) GetUser(ctx context.Context, nip string) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, password, access, role FROM users WHERE nip = $1"

	row := r.db.QueryRow(ctx, query, nip)
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Access, &user.Role)
	if err != nil {
		return nil, err
	}
	user.Nip = nip

	return &user, nil
}

func (r *userRepositories) CreateUser(ctx context.Context, user *models.AdminRegistrationPayload, userId uuid.UUID, refreshToken string) (string, error) {
	var id string
	role := CheckRoleForRegister(strconv.FormatInt(user.Nip, 10))
	statement := "INSERT INTO users (id, name, nip, password, role, refresh_token) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	row := r.db.QueryRow(ctx, statement, userId, user.Name, strconv.FormatInt(user.Nip, 10), user.Password, role, refreshToken)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepositories) CreateUserTx(ctx context.Context, tx pgx.Tx, user *models.AdminRegistrationPayload, hashPassword string) (string, error) {
	var id string
	statement := "INSERT INTO users (name, nip, password) VALUES ($1, $2, $3) RETURNING id"

	row := tx.QueryRow(ctx, statement, user.Name, strconv.FormatInt(user.Nip, 10), hashPassword)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepositories) UpdateRefreshToken(ctx context.Context, userId string, refreshToken string) (pgconn.CommandTag, error) {
	query := "UPDATE users SET refresh_token = $1 WHERE id = $2"

	res, err := r.db.Exec(ctx, query, refreshToken, userId)

	return res, err
}

func (r *userRepositories) GetUserById(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := "SELECT id,nip FROM users WHERE id = $1"

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&user.ID, &user.Nip)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepositories) GetNurseAccessByNip(ctx context.Context, nip string) (*models.User, error) {
	var user models.User
	query := "SELECT access FROM USERS where nip = $1"

	row := r.db.QueryRow(ctx, query, nip)
	err := row.Scan(&user.Access)
	if err != nil {
		return nil, err
	}
	user.Nip = nip

	return &user, nil
}

func CheckRoleForRegister(nip string) string {
	if strings.HasPrefix(nip, "615") {
		return "admin"
	} else if strings.HasPrefix(nip, "303") {
		return "nurse"
	}
	return ""
}
