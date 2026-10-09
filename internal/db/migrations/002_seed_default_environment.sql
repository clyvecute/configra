-- Up
INSERT INTO users (id, email, password_hash)
VALUES (1, 'default@configra.local', 'hash')
ON CONFLICT (id) DO NOTHING;

INSERT INTO projects (id, name, owner_id, api_key)
VALUES (1, 'Default Project', 1, 'valid-key')
ON CONFLICT (id) DO NOTHING;

INSERT INTO environments (id, project_id, name, slug)
VALUES (1, 1, 'Production', 'prod')
ON CONFLICT (id) DO NOTHING;

-- Down
DELETE FROM environments WHERE id = 1;
DELETE FROM projects WHERE id = 1;
DELETE FROM users WHERE id = 1;
