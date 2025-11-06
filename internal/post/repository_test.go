package post

import (
    "regexp"
    "testing"

    sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_DeleteTx_SoftDelete(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil { t.Fatalf("sqlmock.New: %v", err) }
    defer db.Close()

    repo := NewRepository(db)

    mock.ExpectBegin()
    mock.ExpectExec(regexp.QuoteMeta("UPDATE posts SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL")).
        WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
    mock.ExpectCommit()

    tx, _ := db.Begin()
    if err := repo.DeleteTx(tx, 1); err != nil { t.Fatalf("DeleteTx error: %v", err) }
    if err := tx.Commit(); err != nil { t.Fatalf("commit: %v", err) }
    if err := mock.ExpectationsWereMet(); err != nil { t.Fatalf("unmet: %v", err) }
}

func TestRepository_UpdateTx_SoftUpdate(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil { t.Fatalf("sqlmock.New: %v", err) }
    defer db.Close()

    repo := NewRepository(db)

    mock.ExpectBegin()
    // CTE query - we just match INSERT INTO posts ... RETURNING id and provide a row
    mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO posts (user_id, title, content) SELECT user_id, COALESCE($2, title), COALESCE($3, content) FROM old RETURNING id")).
        WithArgs(5, sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
    mock.ExpectCommit()

    tx, _ := db.Begin()
    title := "new title"
    newID, err := repo.UpdateTx(tx, 5, &title, nil)
    if err != nil { t.Fatalf("UpdateTx error: %v", err) }
    if newID != 10 { t.Fatalf("expected newID 10, got %d", newID) }
    if err := tx.Commit(); err != nil { t.Fatalf("commit: %v", err) }
    if err := mock.ExpectationsWereMet(); err != nil { t.Fatalf("unmet: %v", err) }
}


