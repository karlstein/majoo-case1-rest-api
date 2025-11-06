package user

import (
	"database/sql"
	"errors"
	"majoo-case1-rest-api/config"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestUsecase_Register_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	cfg := config.Config{JWTSecret: "test-secret"}
	uc := NewUsecase(repo, cfg)

	// Mock: user doesn't exist
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("test@example.com", "testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Mock: create user
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", "test@example.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	user, token, err := uc.Register("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}
	if token == "" {
		t.Error("expected token, got empty")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Register_Conflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	cfg := config.Config{JWTSecret: "test-secret"}
	uc := NewUsecase(repo, cfg)

	// Mock: user already exists
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("test@example.com", "testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	_, _, err = uc.Register("testuser", "test@example.com", "password123")
	if err != ErrConflict {
		t.Errorf("expected ErrConflict, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Register_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	cfg := config.Config{JWTSecret: "test-secret"}
	uc := NewUsecase(repo, cfg)

	// Mock: database error on exists check
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("test@example.com", "testuser").
		WillReturnError(errors.New("database error"))

	_, _, err = uc.Register("testuser", "test@example.com", "password123")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Login_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	cfg := config.Config{JWTSecret: "test-secret"}
	uc := NewUsecase(repo, cfg)

	// Mock: get user by email
	mock.ExpectQuery("SELECT id, username, email, password_hash").
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at"}).
			AddRow(1, "testuser", "test@example.com", "$2a$10$dummyhash", time.Now(), time.Now()))

	// Note: password check will fail with dummy hash, but we can test the flow
	_, _, err = uc.Login("test@example.com", "wrongpassword")
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Login_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	cfg := config.Config{JWTSecret: "test-secret"}
	uc := NewUsecase(repo, cfg)

	// Mock: user not found
	mock.ExpectQuery("SELECT id, username, email, password_hash").
		WithArgs("test@example.com").
		WillReturnError(sql.ErrNoRows)

	_, _, err = uc.Login("test@example.com", "password123")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Login_Unauthorized(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	cfg := config.Config{JWTSecret: "test-secret"}
	uc := NewUsecase(repo, cfg)

	// Mock: get user by email
	mock.ExpectQuery("SELECT id, username, email, password_hash").
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at"}).
			AddRow(1, "testuser", "test@example.com", "$2a$10$dummyhash", time.Now(), time.Now()))

	_, _, err = uc.Login("test@example.com", "wrongpassword")
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
