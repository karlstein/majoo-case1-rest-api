package post

import "database/sql"

type Usecase struct {
	repo *Repository
	db   *sql.DB
}

func NewUsecase(db *sql.DB, repo *Repository) *Usecase { return &Usecase{db: db, repo: repo} }

func (u *Usecase) List(page, limit int) ([]Post, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	rows, err := u.repo.List(limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt, &p.Author); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (u *Usecase) Get(id int) (Post, error) {
	row, _ := u.repo.GetByID(id)
	var p Post
	if err := row.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt, &p.Author); err != nil {
		return Post{}, err
	}
	return p, nil
}

func (u *Usecase) Create(userID int, req CreatePostRequest) (Post, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return Post{}, err
	}
	defer tx.Rollback()
	id, err := u.repo.CreateTx(tx, userID, req.Title, req.Content)
	if err != nil {
		return Post{}, err
	}
	if err := tx.Commit(); err != nil {
		return Post{}, err
	}
	return u.Get(id)
}

func (u *Usecase) Update(userID, id int, req UpdatePostRequest) (Post, error) {
	ownerID, err := u.repo.GetOwnerID(id)
	if err != nil {
		return Post{}, err
	}
	if ownerID != userID {
		return Post{}, ErrForbidden
	}
	tx, err := u.db.Begin()
	if err != nil {
		return Post{}, err
	}
	defer tx.Rollback()
	newID, err := u.repo.UpdateTx(tx, id, req.Title, req.Content)
	if err != nil {
		return Post{}, err
	}
	if err := tx.Commit(); err != nil {
		return Post{}, err
	}
	return u.Get(newID)
}

func (u *Usecase) Delete(userID, id int) error {
	ownerID, err := u.repo.GetOwnerID(id)
	if err != nil {
		return err
	}
	if ownerID != userID {
		return ErrForbidden
	}
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := u.repo.DeleteTx(tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

var (
	ErrForbidden = errString("forbidden")
)

type errString string

func (e errString) Error() string { return string(e) }
