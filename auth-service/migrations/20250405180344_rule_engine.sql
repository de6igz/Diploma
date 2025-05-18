-- +goose Up
-- +goose StatementBegin

CREATE schema if not exists rule_engine;
CREATE TABLE  if not exists rule_engine.projects (
                                      id          SERIAL PRIMARY KEY,
                                      project_name VARCHAR(255),
                                      description TEXT,
                                      user_id     INTEGER
);

CREATE TABLE if not exists rule_engine.services (
                                      id          SERIAL PRIMARY KEY,
                                      project_id  INTEGER NOT NULL,        -- Ссылка на projects.id
                                      service_name VARCHAR(255),
                                      UNIQUE (project_id, service_name),
                                      FOREIGN KEY (project_id) REFERENCES rule_engine.projects(id)
);

CREATE TABLE if not exists rule_engine.error_rules (
                                         id          SERIAL PRIMARY KEY,
                                         name        VARCHAR(255),
                                         actions     JSONB,
                                         root_node   JSONB,
                                         user_id     INTEGER,
                                         service_id  INTEGER,
                                        description VARCHAR(255),
                                         FOREIGN KEY (service_id) REFERENCES rule_engine.services(id)
);

CREATE TABLE if not exists rule_engine.resource_rules (
                                            id          SERIAL PRIMARY KEY,
                                            name        VARCHAR(255),
                                            actions     JSONB,
                                            root_node   JSONB,
                                            user_id     INTEGER,
                                            service_id  INTEGER,
                                            description VARCHAR(255),
                                            FOREIGN KEY (service_id) REFERENCES rule_engine.services(id)
);




--
-- INSERT INTO rule_engine.projects (project_name, description, user_id)
-- VALUES ('marketplace', 'Публичное API для сервиса уведомлений Apple', 1);
--
-- INSERT INTO rule_engine.services (project_id, service_name)
-- VALUES (1, 'Example Service');
--
--
-- INSERT INTO rule_engine.error_rules (name, actions, root_node, user_id, service_id)
-- VALUES (
--            'Discord for Dev test 2',
--            '[
--              {"type": "TELEGRAM", "params": {"key": "chat_id", "value": "402379231"}},
--              {"type": "TELEGRAM", "params": {"key": "chat_id", "value": "456"}},
--              {"type": "DISCORD", "params": {"key": "server", "value": "789"}}
--            ]'::jsonb,
--            '{
--              "operator": "OR",
--              "conditions": [
--                {"field": "environment", "operator": "eq", "value": "dev"}
--              ],
--              "children": [
--                {
--                  "operator": "OR",
--                  "conditions": [
--                    {"field": "error_message", "operator": "eq", "value": "ERROR_PANIC"}
--                  ],
--                  "children": []
--                }
--              ]
--            }'::jsonb,
--            1,
--            1
--        );
--
-- INSERT INTO rule_engine.error_rules (name, actions, root_node, user_id, service_id)
-- VALUES (
--            'Discord for Dev test 3',
--            '[
--              {"type": "TELEGRAM", "params": {"key": "chat_id", "value": "402379231"}},
--              {"type": "TELEGRAM", "params": {"key": "chat_id", "value": "456"}},
--              {"type": "DISCORD", "params": {"key": "server", "value": "789"}}
--            ]'::jsonb,
--            '{
--              "operator": "OR",
--              "conditions": [
--                {"field": "environment", "operator": "eq", "value": "dev"}
--              ],
--              "children": [
--                {
--                  "operator": "OR",
--                  "conditions": [
--                    {"field": "error_message", "operator": "eq", "value": "ERROR_PANIC"}
--                  ],
--                  "children": []
--                }
--              ]
--            }'::jsonb,
--            1,
--            1
--        );
--
-- INSERT INTO rule_engine.resource_rules (name, actions, root_node, user_id, service_id)
-- VALUES (
--            'Resource Test',
--            '[
--              {"type": "TELEGRAM", "params": {"key": "chat_id", "value": "402379231"}},
--              {"type": "TELEGRAM", "params": {"key": "chat_id", "value": "456"}},
--              {"type": "DISCORD", "params": {"key": "server", "value": "789"}}
--            ]'::jsonb,
--            '{
--              "operator": "OR",
--              "conditions": [
--                {"field": "fields.memory_alloc_bytes", "operator": "gte", "value": 100000},
--                {"field": "fields.goroutine_count", "operator": "gte", "value": 100}
--              ],
--              "children": []
--            }'::jsonb,
--            1,
--            1
--        );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rule_engine.resource_rules;
DROP TABLE IF EXISTS rule_engine.error_rules;
DROP TABLE IF EXISTS rule_engine.services;
DROP TABLE IF EXISTS rule_engine.projects;
-- +goose StatementEnd
