package comment

import (
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestUsecase_Create_PostNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: post doesn't exist
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	_, err = uc.Create(1, 1, "comment")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Create_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: database error on post exists check
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	_, err = uc.Create(1, 1, "comment")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUsecase_Update_Forbidden(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: get owner ID - different user
	mock.ExpectQuery("SELECT user_id FROM comments").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(999))

	content := "updated"
	_, err = uc.Update(1, 1, &content)
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

	// Mock: comment not found
	mock.ExpectQuery("SELECT user_id FROM comments").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	content := "updated"
	_, err = uc.Update(1, 1, &content)
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
	mock.ExpectQuery("SELECT user_id FROM comments").
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

	// Mock: comment not found
	mock.ExpectQuery("SELECT user_id FROM comments").
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

func TestUsecase_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: comment not found
	mock.ExpectQuery("SELECT c.id, c.post_id").
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

func TestUsecase_Create_TransactionError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	uc := NewUsecase(db, repo)

	// Mock: post exists
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Mock: transaction begin fails
	mock.ExpectBegin().WillReturnError(errors.New("tx begin error"))

	_, err = uc.Create(1, 1, "comment")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

