ALTER TABLE posts ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE comments ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts(deleted_at);
CREATE INDEX IF NOT EXISTS idx_comments_deleted_at ON comments(deleted_at);


