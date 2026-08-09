CREATE INDEX IF NOT EXISTS idx_bookmarks_user_created_at
    ON bookmarks (user_id, created_at DESC);
