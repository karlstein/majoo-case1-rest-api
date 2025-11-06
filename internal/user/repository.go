package user

import "database/sql"

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) ExistsByEmailOrUsername(email, username string) (bool, error) {
    var exists bool
    err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)", email, username).Scan(&exists)
    return exists, err
}

func (r *Repository) Create(username, email, passwordHash string) (int, error) {
    var id int
    err := r.db.QueryRow("INSERT INTO users (username, email, password_hash) VALUES ($1,$2,$3) RETURNING id", username, email, passwordHash).Scan(&id)
    return id, err
}

func (r *Repository) GetByEmail(email string) (User, error) {
    var u User
    err := r.db.QueryRow("SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1", email).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
    return u, err
}


