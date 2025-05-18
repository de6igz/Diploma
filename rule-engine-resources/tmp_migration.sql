
CREATE TABLE logs (
                      id SERIAL NOT NULL,
                      user_id INT NOT NULL,
                      service_name VARCHAR NOT NULL,
                      timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
                      log JSONB NOT NULL,
                      PRIMARY KEY (id, timestamp) -- Добавляем timestamp в PRIMARY KEY
);


-- Превращаем в Hypertable (партиционируем по timestamp)
SELECT create_hypertable('logs', 'timestamp');

CREATE INDEX logs_user_idx ON logs (user_id, timestamp DESC);

CREATE INDEX logs_service_idx ON logs (service_name, timestamp DESC);
