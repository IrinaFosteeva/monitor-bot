ALTER TABLE subscriptions
    ADD CONSTRAINT unique_user_target UNIQUE (user_id, target_id);
