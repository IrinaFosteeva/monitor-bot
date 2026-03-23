ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS chat_id BIGINT;

UPDATE subscriptions s
SET chat_id = u.telegram_chat_id
FROM users u
WHERE s.user_id = u.id AND s.chat_id IS NULL;

ALTER TABLE subscriptions ALTER COLUMN chat_id SET NOT NULL;
