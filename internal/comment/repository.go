package comment

import "database/sql"

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) PostExists(id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id=$1 AND deleted_at IS NULL)", id).Scan(&exists)
	return exists, err
}

func (r *Repository) ListByPost(postID int) (*sql.Rows, error) {
	const q = `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, u.username as author
               FROM comments c JOIN users u ON c.user_id = u.id
               WHERE c.post_id = $1 AND c.deleted_at IS NULL
               ORDER BY c.created_at ASC`
	return r.db.Query(q, postID)
}

func (r *Repository) GetByID(id int) (*sql.Row, error) {
	const q = `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, u.username as author
               FROM comments c JOIN users u ON c.user_id = u.id WHERE c.id=$1 AND c.deleted_at IS NULL`
	return r.db.QueryRow(q, id), nil
}

func (r *Repository) GetOwnerID(id int) (int, error) {
	var uid int
	err := r.db.QueryRow("SELECT user_id FROM comments WHERE id=$1 AND deleted_at IS NULL", id).Scan(&uid)
	return uid, err
}

func (r *Repository) CreateTx(tx *sql.Tx, postID, userID int, content string) (int, error) {
	var id int
	err := tx.QueryRow("INSERT INTO comments (post_id, user_id, content) VALUES ($1,$2,$3) RETURNING id", postID, userID, content).Scan(&id)
	return id, err
}

// UpdateTx performs a soft update by invalidating the current row and inserting a new one.
func (r *Repository) UpdateTx(tx *sql.Tx, id int, content *string) (int, error) {
	const q = `WITH old AS (
                    SELECT id, post_id, user_id, content FROM comments WHERE id = $1 AND deleted_at IS NULL FOR UPDATE
                ), del AS (
                    UPDATE comments SET deleted_at = CURRENT_TIMESTAMP WHERE id IN (SELECT id FROM old)
                    RETURNING 1
                )
                INSERT INTO comments (post_id, user_id, content)
                SELECT post_id, user_id, COALESCE($2, content) FROM old
                RETURNING id`
	var newID int
	if err := tx.QueryRow(q, id, content).Scan(&newID); err != nil {
		return 0, err
	}
	return newID, nil
}

func (r *Repository) DeleteTx(tx *sql.Tx, id int) error {
	_, err := tx.Exec("UPDATE comments SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL", id)
	return err
}
