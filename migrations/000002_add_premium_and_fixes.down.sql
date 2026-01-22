DROP INDEX IF EXISTS idx_subscriptions_user_id;

ALTER TABLE targets DROP CONSTRAINT IF EXISTS unique_target_url;

-- Восстанавливаем created_by при откате
ALTER TABLE targets ADD COLUMN IF NOT EXISTS created_by INT REFERENCES users(id);

ALTER TABLE subscriptions DROP COLUMN IF EXISTS chat_id;

ALTER TABLE users DROP COLUMN IF EXISTS is_premium;