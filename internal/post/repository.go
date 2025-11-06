package post

import "database/sql"

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(limit, offset int) (*sql.Rows, error) {
    const q = `SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at, u.username as author
               FROM posts p JOIN users u ON p.user_id = u.id
               WHERE p.deleted_at IS NULL
               ORDER BY p.created_at DESC LIMIT $1 OFFSET $2`
    return r.db.Query(q, limit, offset)
}

func (r *Repository) GetByID(id int) (*sql.Row, error) {
    const q = `SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at, u.username as author
               FROM posts p JOIN users u ON p.user_id = u.id WHERE p.id = $1 AND p.deleted_at IS NULL`
    return r.db.QueryRow(q, id), nil
}

func (r *Repository) GetOwnerID(id int) (int, error) {
    var userID int
    err := r.db.QueryRow("SELECT user_id FROM posts WHERE id = $1 AND deleted_at IS NULL", id).Scan(&userID)
    return userID, err
}

func (r *Repository) CreateTx(tx *sql.Tx, userID int, title, content string) (int, error) {
    var id int
    err := tx.QueryRow("INSERT INTO posts (user_id, title, content) VALUES ($1,$2,$3) RETURNING id", userID, title, content).Scan(&id)
    return id, err
}

// UpdateTx performs a soft update: invalidate current row and create a new one.
func (r *Repository) UpdateTx(tx *sql.Tx, id int, title *string, content *string) (int, error) {
    // Use a CTE to fetch old, soft-delete, and insert a new row with merged values
    const q = `WITH old AS (
                    SELECT id, user_id, title, content FROM posts WHERE id = $1 AND deleted_at IS NULL FOR UPDATE
                ), del AS (
                    UPDATE posts SET deleted_at = CURRENT_TIMESTAMP WHERE id IN (SELECT id FROM old)
                    RETURNING 1
                )
                INSERT INTO posts (user_id, title, content)
                SELECT user_id, COALESCE($2, title), COALESCE($3, content) FROM old
                RETURNING id`
    var newID int
    if err := tx.QueryRow(q, id, title, content).Scan(&newID); err != nil {
        return 0, err
    }
    return newID, nil
}

func (r *Repository) DeleteTx(tx *sql.Tx, id int) error {
    _, err := tx.Exec("UPDATE posts SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL", id)
    return err
}


