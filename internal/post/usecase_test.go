package post

import (
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestUsecase_Update_Forbidden(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: get owner ID - different user
	mock.ExpectQuery("SELECT user_id FROM posts").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(999))

	req := UpdatePostRequest{Title: stringPtr("New Title")}
	_, err = uc.Update(1, 1, req)
	if err != ErrForbidden {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: post not found
	mock.ExpectQuery("SELECT user_id FROM posts").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	req := UpdatePostRequest{Title: stringPtr("New Title")}
	_, err = uc.Update(1, 1, req)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Delete_Forbidden(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: get owner ID - different user
	mock.ExpectQuery("SELECT user_id FROM posts").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(999))

	err = uc.Delete(1, 1)
	if err != ErrForbidden {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: post not found
	mock.ExpectQuery("SELECT user_id FROM posts").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	err = uc.Delete(1, 1)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Create_TransactionError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: transaction begin fails
	mock.ExpectBegin().WillReturnError(errors.New("tx begin error"))

	req := CreatePostRequest{Title: "Test", Content: "Content"}
	_, err = uc.Create(1, req)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: post not found
	mock.ExpectQuery("SELECT p.id, p.user_id").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err = uc.Get(1)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func stringPtr(s string) *string {
	return &s
}

