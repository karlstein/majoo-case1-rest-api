package comment

import "database/sql"

type Usecase struct {
    repo *Repository
    db   *sql.DB
}

func NewUsecase(db *sql.DB, repo *Repository) *Usecase { return &Usecase{db: db, repo: repo} }

func (u *Usecase) ListByPost(postID int) ([]Comment, error) {
    rows, err := u.repo.ListByPost(postID)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Comment
    for rows.Next() {
        var c Comment
        if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.Author); err != nil { return nil, err }
        out = append(out, c)
    }
    return out, nil
}

func (u *Usecase) Get(id int) (Comment, error) {
    row, _ := u.repo.GetByID(id)
    var c Comment
    if err := row.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.Author); err != nil { return Comment{}, err }
    return c, nil
}

func (u *Usecase) Create(postID, userID int, content string) (Comment, error) {
    exists, err := u.repo.PostExists(postID)
    if err != nil { return Comment{}, err }
    if !exists { return Comment{}, ErrNotFound }
    tx, err := u.db.Begin()
    if err != nil { return Comment{}, err }
    defer tx.Rollback()
    id, err := u.repo.CreateTx(tx, postID, userID, content)
    if err != nil { return Comment{}, err }
    if err := tx.Commit(); err != nil { return Comment{}, err }
    return u.Get(id)
}

func (u *Usecase) Update(userID, id int, content *string) (Comment, error) {
    ownerID, err := u.repo.GetOwnerID(id)
    if err != nil { return Comment{}, err }
    if ownerID != userID { return Comment{}, ErrForbidden }
    tx, err := u.db.Begin()
    if err != nil { return Comment{}, err }
    defer tx.Rollback()
    newID, err := u.repo.UpdateTx(tx, id, content)
    if err != nil { return Comment{}, err }
    if err := tx.Commit(); err != nil { return Comment{}, err }
    return u.Get(newID)
}

func (u *Usecase) Delete(userID, id int) error {
    ownerID, err := u.repo.GetOwnerID(id)
    if err != nil { return err }
    if ownerID != userID { return ErrForbidden }
    tx, err := u.db.Begin()
    if err != nil { return err }
    defer tx.Rollback()
    if err := u.repo.DeleteTx(tx, id); err != nil { return err }
    return tx.Commit()
}

var (
    ErrForbidden = errString("forbidden")
    ErrNotFound  = errString("not_found")
)

type errString string
func (e errString) Error() string { return string(e) }


