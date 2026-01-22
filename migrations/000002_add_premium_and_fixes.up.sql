ALTER TABLE users ADD COLUMN IF NOT EXISTS is_premium BOOLEAN DEFAULT FALSE;

ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS chat_id BIGINT;

UPDATE subscriptions s
SET chat_id = u.telegram_chat_id
FROM users u
WHERE s.user_id = u.id AND s.chat_id IS NULL;

ALTER TABLE subscriptions ALTER COLUMN chat_id SET NOT NULL;

DELETE FROM targets a USING targets b
WHERE a.id > b.id AND a.url = b.url;

ALTER TABLE targets ADD CONSTRAINT unique_target_url UNIQUE (url);

-- Удаляем created_by из targets (shared model - targets ничьи!)
ALTER TABLE targets DROP COLUMN IF EXISTS created_by;

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);