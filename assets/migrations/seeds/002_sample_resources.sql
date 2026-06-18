INSERT INTO resources (name, description, status, created_by_id)
SELECT 'Example Resource', 'A generic resource used as a CRUD implementation reference.', 'active', users.id
FROM users
WHERE users.email = 'admin@example.com'
AND NOT EXISTS (
    SELECT 1 FROM resources WHERE resources.name = 'Example Resource'
);

INSERT INTO resources (name, description, status, created_by_id)
SELECT 'Archived Resource', 'A second sample row for filtering and update examples.', 'archived', users.id
FROM users
WHERE users.email = 'admin@example.com'
AND NOT EXISTS (
    SELECT 1 FROM resources WHERE resources.name = 'Archived Resource'
);
