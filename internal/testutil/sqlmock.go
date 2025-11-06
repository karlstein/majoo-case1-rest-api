package testutil

import (
    "database/sql"

    sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func NewSQLMock() (*sql.DB, sqlmock.Sqlmock, error) {
    db, mock, err := sqlmock.New()
    if err != nil {
        return nil, nil, err
    }
    return db, mock, nil
}


