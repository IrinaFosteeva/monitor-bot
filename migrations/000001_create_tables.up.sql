CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       telegram_chat_id BIGINT NOT NULL UNIQUE,
                       is_active BOOLEAN DEFAULT TRUE,
                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE targets (
                         id SERIAL PRIMARY KEY,
                         name TEXT NOT NULL,
                         url TEXT NOT NULL,
                         method TEXT DEFAULT 'GET',
                         expected_status INT DEFAULT 200,
                         body_regex TEXT,
                         interval_seconds INT DEFAULT 60,
                         timeout_seconds INT DEFAULT 5,
                         region_restriction TEXT,
                         created_by INT REFERENCES users(id),
                         enabled BOOLEAN DEFAULT TRUE,
                         created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE checks (
                        id SERIAL PRIMARY KEY,
                        target_id INT REFERENCES targets(id),
                        timestamp TIMESTAMP DEFAULT NOW(),
                        status TEXT,
                        http_code INT,
                        response_time_ms INT,
                        error TEXT,
                        region TEXT
);

CREATE TABLE subscriptions (
                               id SERIAL PRIMARY KEY,
                               user_id INT REFERENCES users(id),
                               target_id INT REFERENCES targets(id),
                               notify_down_only BOOLEAN DEFAULT TRUE,
                               min_retries INT DEFAULT 1,
                               created_at TIMESTAMP DEFAULT NOW()
);
