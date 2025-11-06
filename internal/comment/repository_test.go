package comment

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
    mock.ExpectExec(regexp.QuoteMeta("UPDATE comments SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL")).
        WithArgs(2).WillReturnResult(sqlmock.NewResult(0, 1))
    mock.ExpectCommit()

    tx, _ := db.Begin()
    if err := repo.DeleteTx(tx, 2); err != nil { t.Fatalf("DeleteTx error: %v", err) }
    if err := tx.Commit(); err != nil { t.Fatalf("commit: %v", err) }
    if err := mock.ExpectationsWereMet(); err != nil { t.Fatalf("unmet: %v", err) }
}

func TestRepository_UpdateTx_SoftUpdate(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil { t.Fatalf("sqlmock.New: %v", err) }
    defer db.Close()

    repo := NewRepository(db)

    mock.ExpectBegin()
    mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO comments (post_id, user_id, content) SELECT post_id, user_id, COALESCE($2, content) FROM old RETURNING id")).
        WithArgs(9, sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(11))
    mock.ExpectCommit()

    tx, _ := db.Begin()
    content := "updated"
    newID, err := repo.UpdateTx(tx, 9, &content)
    if err != nil { t.Fatalf("UpdateTx error: %v", err) }
    if newID != 11 { t.Fatalf("expected newID 11, got %d", newID) }
    if err := tx.Commit(); err != nil { t.Fatalf("commit: %v", err) }
    if err := mock.ExpectationsWereMet(); err != nil { t.Fatalf("unmet: %v", err) }
}


